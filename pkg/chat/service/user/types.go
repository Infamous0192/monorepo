package user

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/types/pagination"
	"context"
)

// UserService defines the interface for user-related operations
type UserService interface {
	// Get retrieves a user by ID
	Get(ctx context.Context, id string) (*entity.User, error)

	// GetAll retrieves multiple users with pagination
	GetAll(ctx context.Context, pag pagination.Pagination) ([]*entity.User, int64, error)

	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error

	// Update modifies an existing user
	Update(ctx context.Context, user *entity.User) error

	// Delete removes a user
	Delete(ctx context.Context, id string) error
}
