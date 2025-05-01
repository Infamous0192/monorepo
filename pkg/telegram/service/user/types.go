package user

import (
	"app/pkg/telegram/domain/entity"
	"app/pkg/types/pagination"
	"context"
)

// UserService defines the interface for Telegram user operations
type UserService interface {
	// Get retrieves a single user by ID
	Get(ctx context.Context, id string) (*entity.User, error)

	// GetByTelegramID retrieves a single user by Telegram ID
	GetByTelegramID(ctx context.Context, telegramID int64) (*entity.User, error)

	// GetByClientID retrieves users by client ID with pagination
	GetByClientID(ctx context.Context, clientID string, pag pagination.Pagination) ([]*entity.User, int64, error)

	// GetAll retrieves multiple users with pagination
	GetAll(ctx context.Context, pag pagination.Pagination) ([]*entity.User, int64, error)

	// Create stores a new user
	Create(ctx context.Context, user *entity.User) error

	// Update modifies an existing user
	Update(ctx context.Context, user *entity.User) error

	// Delete removes a user
	Delete(ctx context.Context, id string) error
}
