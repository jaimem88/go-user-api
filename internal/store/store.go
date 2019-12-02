package store

import (
	"sync"

	user_pb "github.com/AssemblyPayments/go-user-api/pkg/proto"
)

// InMemoryStore implements the Store interface defined in the API
type InMemoryStore struct {
	usersMux *sync.RWMutex
	users    map[string]*user_pb.User
}

func New() *InMemoryStore {
	return &InMemoryStore{
		usersMux: new(sync.RWMutex),
		users:    make(map[string]*user_pb.User),
	}
}

// CreateUser creates a new user_pb.User and stores it in memory
// It returns the created user
func (s *InMemoryStore) CreateUser(newUser *user_pb.CreateUserRequest) (*user_pb.User, error) {
	return nil, nil
}

// GetUser gets a user_pb.User by its ID
func (s *InMemoryStore) GetUser(id string) (*user_pb.User, error) {
	return nil, nil
}
