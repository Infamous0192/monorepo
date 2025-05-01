package chat

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/chat/domain/repository"
	"app/pkg/types/pagination"
	"context"
)

// SendMessageParams represents parameters for sending a message
type SendMessageParams struct {
	ChatroomID string
	SenderID   string
	Message    string
}

// SendDirectMessageParams represents parameters for sending a direct message
type SendDirectMessageParams struct {
	SenderID   string
	ReceiverID string
	Message    string
}

// ChatService defines the interface for chat-related operations
type ChatService interface {
	// GetChat retrieves a single chat message by ID
	GetChat(ctx context.Context, chatID string) (*entity.Chat, error)

	// GetChats retrieves multiple chat messages with filtering and pagination
	GetChats(ctx context.Context, filter repository.ChatFilter, pag pagination.Pagination) ([]*entity.Chat, int64, error)

	// UpdateChat modifies an existing chat message
	UpdateChat(ctx context.Context, chat *entity.Chat) error

	// RemoveChat deletes a chat message by ID
	RemoveChat(ctx context.Context, chatID string) error

	// SendMessage sends a message to a chatroom
	// It will validate:
	// - The chatroom exists
	// - The sender is a participant in the chatroom
	// - The sender is not muted
	// Returns the created chat message
	SendMessage(ctx context.Context, params SendMessageParams) (*entity.Chat, error)

	// SendDirectMessage sends a direct message to another user
	// It will:
	// - Create a direct chatroom if it doesn't exist
	// - Add both users as participants
	// - Send the message
	// Returns the created chat message and the chatroom
	SendDirectMessage(ctx context.Context, params SendDirectMessageParams) (*entity.Chat, *entity.Chatroom, error)
}
