package chatroom

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/chat/domain/repository"
	"app/pkg/types/pagination"
	"context"
	"time"
)

// CreateChatroomParams represents data for creating a new chatroom
type CreateChatroomParams struct {
	Name         string
	Type         entity.ChatroomType
	IsGroup      bool
	Creator      string   // Creator's participant ID
	Participants []string // List of participant IDs
}

// UpdateChatroomParams represents data for updating a chatroom
type UpdateChatroomParams struct {
	ID   string
	Name *string
	Type *entity.ChatroomType
}

// MuteParticipantParams represents data for muting a participant
type MuteParticipantParams struct {
	ChatroomID    string
	ParticipantID string
	Duration      time.Duration // How long to mute the participant
}

// UpdateParticipantRoleParams represents data for updating a participant's role
type UpdateParticipantRoleParams struct {
	ChatroomID    string
	ParticipantID string
	NewRole       entity.ParticipantRole
}

// ChatroomService defines the interface for chatroom-related operations
type ChatroomService interface {
	// GetChatroom retrieves a single chatroom by ID with populated participant references
	GetChatroom(ctx context.Context, id string) (*entity.ChatroomPopulated, error)

	// GetChatrooms retrieves multiple chatrooms with filtering and pagination
	GetChatrooms(ctx context.Context, filter repository.ChatroomFilter, pag pagination.Pagination) ([]*entity.ChatroomPopulated, int64, error)

	// CreateChatroom creates a new chatroom
	// It will:
	// - Create the chatroom with the given parameters
	// - Add the creator as an admin participant
	// - Add all specified participants as regular participants
	CreateChatroom(ctx context.Context, data CreateChatroomParams) (*entity.Chatroom, error)

	// UpdateChatroom modifies an existing chatroom
	// Only admin and super_admin participants can update chatroom properties
	UpdateChatroom(ctx context.Context, data UpdateChatroomParams) (*entity.Chatroom, error)

	// DeleteChatroom removes a chatroom and all its messages
	// Only super_admin participants can delete a chatroom
	DeleteChatroom(ctx context.Context, id string) error

	// AddParticipant adds a participant to a chatroom
	// It will validate:
	// - The chatroom exists
	// - The participant exists
	// - The participant is not already in the chatroom
	AddParticipant(ctx context.Context, chatroomID string, participantID string) error

	// RemoveParticipant removes a participant from a chatroom
	// It will validate:
	// - The chatroom exists
	// - The participant is in the chatroom
	// - For group chats, admin/super_admin can remove anyone with lower roles, regular participants can only remove themselves
	// - For direct chats, either participant can leave
	RemoveParticipant(ctx context.Context, chatroomID string, participantID string) error

	// UpdateParticipantRole updates a participant's role in the chatroom
	// Only super_admin can update roles, and they can't:
	// - Change their own role
	// - Promote someone to a role higher than their own
	// - Modify someone with a role higher than their own
	UpdateParticipantRole(ctx context.Context, data UpdateParticipantRoleParams) error

	// MuteParticipant temporarily mutes a participant
	// Only admin and super_admin can mute participants with lower roles
	MuteParticipant(ctx context.Context, data MuteParticipantParams) error

	// UnmuteParticipant removes a mute from a participant
	// Only admin and super_admin can unmute participants with lower roles
	UnmuteParticipant(ctx context.Context, chatroomID string, participantID string) error
}
