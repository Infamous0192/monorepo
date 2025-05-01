package repository

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/types/pagination"
	"context"
)

type UserRepository interface {
	// Get retrieves a single user message by ID
	Get(ctx context.Context, id string) (*entity.User, error)

	// GetAll retrieves multiple user messages with pagination
	GetAll(ctx context.Context, pagination pagination.Pagination) ([]*entity.User, int64, error)

	// Create stores a new user message
	Create(ctx context.Context, user *entity.User) error

	// Update modifies an existing user message
	Update(ctx context.Context, user *entity.User) error

	// Delete removes a user message
	Delete(ctx context.Context, id string) error
}
