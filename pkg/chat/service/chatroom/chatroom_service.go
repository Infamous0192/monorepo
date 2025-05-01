package chatroom

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/chat/domain/repository"
	"app/pkg/types/pagination"
	"context"
	"fmt"
	"time"
)

type chatroomService struct {
	chatroomRepo repository.ChatroomRepository
}

// NewChatroomService creates a new instance of ChatroomService
func NewChatroomService(chatroomRepo repository.ChatroomRepository) ChatroomService {
	return &chatroomService{
		chatroomRepo: chatroomRepo,
	}
}

// GetChatroom retrieves a single chatroom by ID with populated participant references
func (s *chatroomService) GetChatroom(ctx context.Context, id string) (*entity.ChatroomPopulated, error) {
	return s.chatroomRepo.GetPopulated(ctx, id)
}

// GetChatrooms retrieves multiple chatrooms with filtering and pagination
func (s *chatroomService) GetChatrooms(ctx context.Context, filter repository.ChatroomFilter, pag pagination.Pagination) ([]*entity.ChatroomPopulated, int64, error) {
	return s.chatroomRepo.GetAllPopulated(ctx, filter, pag)
}

// CreateChatroom creates a new chatroom
func (s *chatroomService) CreateChatroom(ctx context.Context, data CreateChatroomParams) (*entity.Chatroom, error) {
	// Create the chatroom
	newChatroom := &entity.Chatroom{
		Name:             data.Name,
		Type:             data.Type,
		IsGroup:          data.IsGroup,
		CreatedTimestamp: time.Now().Unix(),
		MessagesCount:    0,
		Participants:     make([]entity.ChatroomParticipant, 0, len(data.Participants)),
	}

	// Add creator as admin
	creatorParticipant := entity.ChatroomParticipant{
		User:            data.Creator,
		Role:            entity.ParticipantRoleAdmin,
		JoinedTimestamp: time.Now().Unix(),
	}
	newChatroom.Participants = append(newChatroom.Participants, creatorParticipant)

	// Add other participants as members
	for _, participantID := range data.Participants {
		if participantID != data.Creator {
			participant := entity.ChatroomParticipant{
				User:            participantID,
				Role:            entity.ParticipantRoleMember,
				JoinedTimestamp: time.Now().Unix(),
			}
			newChatroom.Participants = append(newChatroom.Participants, participant)
		}
	}

	// Save the chatroom
	if err := s.chatroomRepo.Create(ctx, newChatroom); err != nil {
		return nil, err
	}

	return newChatroom, nil
}

// UpdateChatroom modifies an existing chatroom
func (s *chatroomService) UpdateChatroom(ctx context.Context, data UpdateChatroomParams) (*entity.Chatroom, error) {
	// Get existing chatroom
	chatroom, err := s.chatroomRepo.Get(ctx, data.ID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if data.Name != nil {
		chatroom.Name = *data.Name
	}
	if data.Type != nil {
		chatroom.Type = *data.Type
	}

	// Save changes
	if err := s.chatroomRepo.Update(ctx, chatroom); err != nil {
		return nil, err
	}

	return chatroom, nil
}

// DeleteChatroom removes a chatroom and all its messages
func (s *chatroomService) DeleteChatroom(ctx context.Context, id string) error {
	// Get chatroom to validate it exists
	chatroom, err := s.chatroomRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	// Delete the chatroom
	return s.chatroomRepo.Delete(ctx, chatroom.ID)
}

// AddParticipant adds a participant to a chatroom
func (s *chatroomService) AddParticipant(ctx context.Context, chatroomID string, participantID string) error {
	// Get chatroom to validate it exists
	chatroom, err := s.chatroomRepo.Get(ctx, chatroomID)
	if err != nil {
		return err
	}

	// Check if participant already exists
	for _, p := range chatroom.Participants {
		if p.User == participantID {
			return fmt.Errorf("participant already exists in chatroom")
		}
	}

	// Create new participant
	participant := entity.ChatroomParticipant{
		User:            participantID,
		Role:            entity.ParticipantRoleMember,
		JoinedTimestamp: time.Now().Unix(),
	}

	return s.chatroomRepo.AddParticipant(ctx, chatroomID, participant)
}

// RemoveParticipant removes a participant from a chatroom
func (s *chatroomService) RemoveParticipant(ctx context.Context, chatroomID string, participantID string) error {
	// Get chatroom to validate it exists and check roles
	chatroom, err := s.chatroomRepo.Get(ctx, chatroomID)
	if err != nil {
		return err
	}

	// For direct chats, either participant can leave
	if !chatroom.IsGroup {
		return s.chatroomRepo.RemoveParticipant(ctx, chatroomID, participantID)
	}

	// For group chats, validate roles
	var participantRole entity.ParticipantRole
	var found bool
	for _, p := range chatroom.Participants {
		if p.User == participantID {
			participantRole = p.Role
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("participant not found in chatroom")
	}

	// Regular members can only remove themselves
	if participantRole == entity.ParticipantRoleMember {
		return s.chatroomRepo.RemoveParticipant(ctx, chatroomID, participantID)
	}

	return fmt.Errorf("operation not allowed for participant role")
}

// UpdateParticipantRole updates a participant's role in the chatroom
func (s *chatroomService) UpdateParticipantRole(ctx context.Context, data UpdateParticipantRoleParams) error {
	// Get chatroom to validate it exists and check roles
	chatroom, err := s.chatroomRepo.Get(ctx, data.ChatroomID)
	if err != nil {
		return err
	}

	// Find the participant to update
	var participant *entity.ChatroomParticipant
	for i := range chatroom.Participants {
		if chatroom.Participants[i].User == data.ParticipantID {
			participant = &chatroom.Participants[i]
			break
		}
	}

	if participant == nil {
		return fmt.Errorf("participant not found in chatroom")
	}

	// Update the role
	participant.Role = data.NewRole

	return s.chatroomRepo.UpdateParticipant(ctx, data.ChatroomID, *participant)
}

// MuteParticipant temporarily mutes a participant
func (s *chatroomService) MuteParticipant(ctx context.Context, data MuteParticipantParams) error {
	// Get chatroom to validate it exists and check roles
	chatroom, err := s.chatroomRepo.Get(ctx, data.ChatroomID)
	if err != nil {
		return err
	}

	// Find the participant to mute
	var participant *entity.ChatroomParticipant
	for i := range chatroom.Participants {
		if chatroom.Participants[i].User == data.ParticipantID {
			participant = &chatroom.Participants[i]
			break
		}
	}

	if participant == nil {
		return fmt.Errorf("participant not found in chatroom")
	}

	// Calculate mute end time
	muteUntil := time.Now().Add(data.Duration).Unix()
	participant.MutedUntilTimestamp = &muteUntil

	return s.chatroomRepo.UpdateParticipant(ctx, data.ChatroomID, *participant)
}

// UnmuteParticipant removes a mute from a participant
func (s *chatroomService) UnmuteParticipant(ctx context.Context, chatroomID string, participantID string) error {
	// Get chatroom to validate it exists and check roles
	chatroom, err := s.chatroomRepo.Get(ctx, chatroomID)
	if err != nil {
		return err
	}

	// Find the participant to unmute
	var participant *entity.ChatroomParticipant
	for i := range chatroom.Participants {
		if chatroom.Participants[i].User == participantID {
			participant = &chatroom.Participants[i]
			break
		}
	}

	if participant == nil {
		return fmt.Errorf("participant not found in chatroom")
	}

	// Remove mute
	participant.MutedUntilTimestamp = nil

	return s.chatroomRepo.UpdateParticipant(ctx, chatroomID, *participant)
}
