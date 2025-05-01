package user

import (
	"app/pkg/quiz/domain/entity"
	"app/pkg/types/pagination"
	"context"
)

// UserService defines the interface for user-related operations
type UserService interface {
	// FindOne retrieves a user by ID
	FindOne(ctx context.Context, id uint) (*entity.User, error)

	// FindAll retrieves multiple users with filtering and pagination
	FindAll(ctx context.Context, query entity.UserQuery) (*pagination.PaginatedResult[entity.User], error)

	// Create creates a new user
	Create(ctx context.Context, userDTO entity.UserDTO) (*entity.User, error)

	// Update modifies an existing user
	Update(ctx context.Context, id uint, userDTO entity.UserDTO) (*entity.User, error)

	// Delete removes a user
	Delete(ctx context.Context, id uint) error
}
