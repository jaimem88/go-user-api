package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	user_pb "github.com/jaime/go-user-api/pkg/proto"
)

// Store defines the operations that the API requires to interact with the store to manage users
type Store interface {
	CreateUser(newUser *user_pb.CreateUserRequest) (*user_pb.User, error)
	GetUser(id string) (*user_pb.User, error)
}

// Option to be set in the API
type Option func(*API)

// API impl all rpc calls defined by our proto
type API struct {
	store Store
}

// New - returns a new instance of API
func New(store Store, options ...Option) (*API, error) {

	api := &API{
		store: store,
	}

	for _, o := range options {
		o(api)
	}

	return api, nil
}

// CreateUser implements the User service interface that creates a user and stores it.
// It returns a created user with its ID or an error if something goes wrong
func (a *API) CreateUser(ctx context.Context, req *user_pb.CreateUserRequest) (*user_pb.CreateUserResponse, error) {

	// nolint
	//TODO add implementation
	a.store.CreateUser(req)

	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}

// GetUser implements the User service interface that gets a user from the store.
// It returns a user if found or an error otherwise
func (a *API) GetUser(ctx context.Context, req *user_pb.GetUserRequest) (*user_pb.GetUserResponse, error) {

	// nolint
	//TODO add implementation
	a.store.GetUser(req.Id)

	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
