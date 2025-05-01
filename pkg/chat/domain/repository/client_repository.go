package repository

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/types/pagination"
	"context"
)

// ClientRepository defines the interface for client data access
type ClientRepository interface {
	// Get retrieves a single client by ID
	Get(ctx context.Context, id string) (*entity.Client, error)

	// GetByKey retrieves a single client by client key
	GetByKey(ctx context.Context, clientKey string) (*entity.Client, error)

	// GetAll retrieves multiple clients with pagination
	GetAll(ctx context.Context, pagination pagination.Pagination) ([]*entity.Client, int64, error)

	// Create stores a new client
	Create(ctx context.Context, client *entity.Client) error

	// Update modifies an existing client
	Update(ctx context.Context, client *entity.Client) error

	// Delete removes a client
	Delete(ctx context.Context, id string) error
}
