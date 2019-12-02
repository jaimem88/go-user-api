// +build endtoend

package endtoend

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	user_pb "github.com/AssemblyPayments/go-user-api/pkg/proto"
)

const defaultGRPCClientTimeout = 5 * time.Second

var userClient user_pb.UserServiceClient

func TestMain(m *testing.M) {
	target := os.Getenv("SERVICE_ADDR")
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Errorf("did not connect: %v", err))
	}
	// nolint: errcheck
	defer conn.Close()
	userClient = user_pb.NewUserServiceClient(conn)

	os.Exit(m.Run())
}

func TestCreateUser(t *testing.T) {

	req := user_pb.CreateUserRequest{
		FirstName: "First",
		LastName:  "Last",
		Email:     "test@email.com",
		Admin:     false,
	}
	tcs := []struct {
		name                 string
		req                  *user_pb.CreateUserRequest
		expectedResponse     *user_pb.CreateUserResponse
		expectedStatusCode   codes.Code
		expectedErrorMessage string
	}{
		{
			name:                 "unimplemented",
			req:                  &req,
			expectedResponse:     nil, // on purpose
			expectedStatusCode:   codes.Unimplemented,
			expectedErrorMessage: "method CreateUser not implemented",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), defaultGRPCClientTimeout)
			defer cancel()
			actualResponse, err := userClient.CreateUser(ctx, tc.req)
			if tc.expectedErrorMessage != "" {
				require.NotNil(t, err)
				grpcStatus, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.expectedStatusCode, grpcStatus.Code())
				assert.Equal(t, tc.expectedErrorMessage, grpcStatus.Message())
			} else {
				require.Nil(t, err)
				require.NotNil(t, err)
				diff := deep.Equal(tc.expectedResponse, actualResponse)
				assert.Nil(t, diff)

			}
		})
	}
}

func TestGetUser(t *testing.T) {

	req := user_pb.GetUserRequest{
		Id: "user_id",
	}
	tcs := []struct {
		name                 string
		req                  *user_pb.GetUserRequest
		expectedResponse     *user_pb.GetUserResponse
		expectedStatusCode   codes.Code
		expectedErrorMessage string
	}{
		{
			name:                 "unimplemented",
			req:                  &req,
			expectedResponse:     nil, // on purpose
			expectedStatusCode:   codes.Unimplemented,
			expectedErrorMessage: "method GetUser not implemented",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), defaultGRPCClientTimeout)
			defer cancel()
			actualResponse, err := userClient.GetUser(ctx, tc.req)
			if tc.expectedErrorMessage != "" {
				require.NotNil(t, err)
				grpcStatus, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tc.expectedStatusCode, grpcStatus.Code())
				assert.Equal(t, tc.expectedErrorMessage, grpcStatus.Message())
			} else {
				require.Nil(t, err)
				require.NotNil(t, err)
				diff := deep.Equal(tc.expectedResponse, actualResponse)
				assert.Nil(t, diff)

			}
		})
	}
}
