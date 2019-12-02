package api

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	user_pb "github.com/AssemblyPayments/go-user-api/pkg/proto"
)

func TestAPI_CreateUser(t *testing.T) {
	req := user_pb.CreateUserRequest{
		FirstName: "First",
		LastName:  "Last",
		Email:     "test@email.com",
		Admin:     false,
	}
	tcs := []struct {
		name                 string
		store                *storeMock
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
			a, err := New(tc.store)
			require.NoError(t, err)
			actualResponse, err := a.CreateUser(context.Background(), tc.req)
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

func TestAPI_GetUser(t *testing.T) {
	req := user_pb.GetUserRequest{
		Id: "user_id",
	}
	tcs := []struct {
		name                 string
		store                *storeMock
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
			a, err := New(tc.store)
			require.NoError(t, err)
			actualResponse, err := a.GetUser(context.Background(), tc.req)
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

type storeMock struct{}

func (sm *storeMock) CreateUser(newUser *user_pb.CreateUserRequest) (*user_pb.User, error) {
	return nil, nil
}
func (sm *storeMock) GetUser(id string) (*user_pb.User, error) {
	return nil, nil
}
