package repository

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/types/pagination"
	"context"
)

// ChatroomFilter represents filtering options for chatroom queries
type ChatroomFilter struct {
	ParticipantID string
	Type          *entity.ChatroomType
	IsGroup       *bool
	StartTime     *int64
	EndTime       *int64
}

type ChatroomRepository interface {
	// Get retrieves a single chatroom by ID
	Get(ctx context.Context, id string) (*entity.Chatroom, error)

	// GetPopulated retrieves a single chatroom with populated user references
	GetPopulated(ctx context.Context, id string) (*entity.ChatroomPopulated, error)

	// GetAll retrieves multiple chatrooms with filtering and pagination
	GetAll(ctx context.Context, filter ChatroomFilter, pagination pagination.Pagination) ([]*entity.Chatroom, int64, error)

	// GetAllPopulated retrieves multiple chatrooms with populated user references
	GetAllPopulated(ctx context.Context, filter ChatroomFilter, pagination pagination.Pagination) ([]*entity.ChatroomPopulated, int64, error)

	// Create stores a new chatroom
	Create(ctx context.Context, chatroom *entity.Chatroom) error

	// Update modifies an existing chatroom
	Update(ctx context.Context, chatroom *entity.Chatroom) error

	// Delete removes a chatroom
	Delete(ctx context.Context, id string) error

	// AddParticipant adds a participant to a chatroom
	AddParticipant(ctx context.Context, chatroomID string, participant entity.ChatroomParticipant) error

	// RemoveParticipant removes a participant from a chatroom
	RemoveParticipant(ctx context.Context, chatroomID string, userID string) error

	// UpdateParticipant updates a participant's properties in a chatroom
	UpdateParticipant(ctx context.Context, chatroomID string, participant entity.ChatroomParticipant) error
}
