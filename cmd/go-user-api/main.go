package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/jaime/go-user-api/internal/api"
	"github.com/jaime/go-user-api/internal/env"
	"github.com/jaime/go-user-api/internal/store"
	user_pb "github.com/jaime/go-user-api/pkg/proto"
)

const (
	componentName = "go_user_api"
)

// grpcServerKeepalive is a duration that if the server doesn't see any activity after this time, it pings the client to see if the transport is still alive.
// https://docs.aws.amazon.com/elasticloadbalancing/latest/network/network-load-balancers.html#connection-idle-timeout
// we use AWS Network Load Balancer that has fixed, unmodifiable 350 sec idle timeout.
// We are enabling keep alive to prevent load balancer to reset the connection
// https://translate.google.com.au/translate?hl=en&sl=ja&u=https://hori-ryota.com/blog/failed-to-grpc-with-nlb/&prev=search
// Setting 150 sec, which is less than half.
const grpcServerKeepalive = 150 * time.Second

// const grpcClientKeepalive = 150 * time.Second

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	mlog := logrus.WithFields(logrus.Fields{
		"component": componentName,
		"version":   env.Version(),
	})

	grpc_logrus.ReplaceGrpcLogger(mlog.WithField("component", componentName+"_grpc"))
	mlog.Infof("Starting %s", componentName)

	grpcServer, err := createGRPCServer(mlog)
	if err != nil {
		mlog.WithError(err).Fatal("failed to create grpc server")
	}
	// Start go routines
	go handleExitSignals(grpcServer, mlog)
	serveGRPC(env.ServiceAddr(), grpcServer, mlog)
}

func createGRPCServer(mlog *logrus.Entry) (*grpc.Server, error) {
	interceptors := []grpc.UnaryServerInterceptor{} // order matters

	interceptors = append(interceptors,
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_logrus.UnaryServerInterceptor(mlog), // logs stats on each end point
		grpc_logrus.PayloadUnaryServerInterceptor(mlog, alwaysLoggingDeciderServer),
	)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptors...)),
		grpc.KeepaliveParams(keepalive.ServerParameters{Time: grpcServerKeepalive}),
	)

	store := store.New()

	a, err := api.New(store)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create api")
	}
	user_pb.RegisterUserServiceServer(s, a)

	reflection.Register(s)

	return s, nil
}

func serveGRPC(serviceAddr string, s *grpc.Server, mlog *logrus.Entry) {
	lis, err := net.Listen("tcp", serviceAddr)
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	mlog.Infof("grpc   :%s listening...", serviceAddr)
	err = s.Serve(lis)
	if err != nil {
		mlog.Warnf("gRPC terminating Serve() %s", err)
	}
}

func handleExitSignals(grpcServer *grpc.Server, mlog *logrus.Entry) {
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	<-exitCh

	mlog.Info("gracefully stopping server")
	grpcServer.GracefulStop()
	mlog.Info("server stopped")

}

func alwaysLoggingDeciderServer(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
	return true
}
