package client

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/types/pagination"
	"context"
)

// ClientService defines the interface for client-related operations
type ClientService interface {
	// GetGetClient retrieves a client by ID
	GetClient(ctx context.Context, id string) (*entity.Client, error)

	// GetByKey retrieves a client by client key
	GetByKey(ctx context.Context, clientKey string) (*entity.Client, error)

	// GetClients retrieves multiple clients with pagination
	GetClients(ctx context.Context, pag pagination.Pagination) ([]*entity.Client, int64, error)

	// CreateClient creates a new client
	CreateClient(ctx context.Context, client *entity.Client) error

	// UpdateClient modifies an existing client
	UpdateClient(ctx context.Context, client *entity.Client) error

	// DeleteClient removes a client
	DeleteClient(ctx context.Context, id string) error

	// Authenticate validates a JWT token against the client's auth endpoint and returns the authenticated user
	Authenticate(ctx context.Context, clientID string, token string) (*entity.User, error)

	// ValidateKey validates a client key and returns the client if valid
	ValidateKey(ctx context.Context, clientKey string) (*entity.Client, error)
}
