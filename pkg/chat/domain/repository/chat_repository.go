package repository

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/types/pagination"
	"context"
)

// ChatFilter represents filtering options for chat queries
type ChatFilter struct {
	ChatroomID string
	SenderID   string
	ReceiverID string
	StartTime  *int64
	EndTime    *int64
}

type ChatRepository interface {
	// Get retrieves a single chat message by ID
	Get(ctx context.Context, id string) (*entity.Chat, error)

	// GetPopulated retrieves a single chat message with populated user references
	GetPopulated(ctx context.Context, id string) (*entity.ChatPopulated, error)

	// GetAll retrieves multiple chat messages with filtering and pagination
	GetAll(ctx context.Context, filter ChatFilter, pagination pagination.Pagination) ([]*entity.Chat, int64, error)

	// GetAllPopulated retrieves multiple chat messages with populated user references
	GetAllPopulated(ctx context.Context, filter ChatFilter, pagination pagination.Pagination) ([]*entity.ChatPopulated, int64, error)

	// Create stores a new chat message
	Create(ctx context.Context, chat *entity.Chat) error

	// Update modifies an existing chat message
	Update(ctx context.Context, chat *entity.Chat) error

	// Delete removes a chat message
	Delete(ctx context.Context, id string) error
}
