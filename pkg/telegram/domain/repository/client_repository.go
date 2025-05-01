package repository

import (
	"app/pkg/telegram/domain/entity"
	"app/pkg/types/pagination"
	"context"
)

// ClientRepository defines the interface for Telegram bot client data access
type ClientRepository interface {
	// Get retrieves a single bot client by ID
	Get(ctx context.Context, id string) (*entity.Client, error)

	// GetByToken retrieves a single bot client by token
	GetByToken(ctx context.Context, token string) (*entity.Client, error)

	// GetAll retrieves multiple bot clients with pagination
	GetAll(ctx context.Context, pagination pagination.Pagination) ([]*entity.Client, int64, error)

	// Create stores a new bot client
	Create(ctx context.Context, client *entity.Client) error

	// Update modifies an existing bot client
	Update(ctx context.Context, client *entity.Client) error

	// Delete removes a bot client
	Delete(ctx context.Context, id string) error
}
