// Driving and Driven Ports for User Management
package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/tasjen/fz/api-hexa/internal/domain/model"
)

// Driving Port: UserService defines the methods that the application layer will use to interact with user-related operations.
type UserService interface {
	RegisterUser(ctx context.Context, user model.User) (model.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (model.User, error)
}

// Driven Port: UserRepository defines the methods that the infrastructure layer must implement to persist and retrieve user data.
type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (model.User, error)
}
