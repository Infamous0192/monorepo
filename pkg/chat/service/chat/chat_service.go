package chat

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/chat/domain/repository"
	"app/pkg/chat/service/chatroom"
	"app/pkg/exception"
	"app/pkg/types/pagination"
	"context"
	"fmt"
	"time"
)

type chatService struct {
	chatRepository  repository.ChatRepository
	chatroomService chatroom.ChatroomService
}

// NewChatService creates a new instance of ChatService
func NewChatService(chatRepository repository.ChatRepository, chatroomService chatroom.ChatroomService) ChatService {
	return &chatService{
		chatRepository:  chatRepository,
		chatroomService: chatroomService,
	}
}

// GetChat retrieves a single chat message by ID
func (s *chatService) GetChat(ctx context.Context, chatID string) (*entity.Chat, error) {
	chat, err := s.chatRepository.Get(ctx, chatID)
	if err != nil {
		return nil, exception.NotFound("chatroom")
	}

	return chat, nil
}

// GetChats retrieves multiple chat messages with filtering and pagination
func (s *chatService) GetChats(ctx context.Context, filter repository.ChatFilter, pag pagination.Pagination) ([]*entity.Chat, int64, error) {
	return s.chatRepository.GetAll(ctx, filter, pag)
}

// UpdateChat modifies an existing chat message
func (s *chatService) UpdateChat(ctx context.Context, chat *entity.Chat) error {
	return s.chatRepository.Update(ctx, chat)
}

// RemoveChat deletes a chat message by ID
func (s *chatService) RemoveChat(ctx context.Context, chatID string) error {
	chat, err := s.chatRepository.Get(ctx, chatID)
	if err != nil {
		return exception.NotFound("chatroom")
	}

	return s.chatRepository.Delete(ctx, chat.ID)
}

// SendMessage sends a message to a chatroom
func (s *chatService) SendMessage(ctx context.Context, params SendMessageParams) (*entity.Chat, error) {
	// Validate chatroom exists and user is participant
	chatroom, err := s.chatroomService.GetChatroom(ctx, params.ChatroomID)
	if err != nil {
		return nil, err
	}

	// Validate sender is a participant and not muted
	var isParticipant bool
	for _, p := range chatroom.Participants {
		if p.User.ID == params.SenderID {
			if p.MutedUntilTimestamp != nil {
				mutedUntil := time.Unix(*p.MutedUntilTimestamp, 0)
				if mutedUntil.After(time.Now()) {
					return nil, fmt.Errorf("user is muted")
				}
			}
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, fmt.Errorf("sender is not a participant in the chatroom")
	}

	// Create and save the chat message
	newChat := &entity.Chat{
		Message:          params.Message,
		Sender:           params.SenderID,
		Chatroom:         params.ChatroomID,
		CreatedTimestamp: time.Now().Unix(),
	}

	if err := s.chatRepository.Create(ctx, newChat); err != nil {
		return nil, err
	}

	return newChat, nil
}

// SendDirectMessage sends a direct message to another user
func (s *chatService) SendDirectMessage(ctx context.Context, params SendDirectMessageParams) (*entity.Chat, *entity.Chatroom, error) {
	chatroomType := entity.ChatroomTypePrivate

	// Find or create direct chatroom
	filter := repository.ChatroomFilter{
		Type:          &chatroomType,
		ParticipantID: params.SenderID,
	}

	chatrooms, _, err := s.chatroomService.GetChatrooms(ctx, filter, pagination.Pagination{})
	if err != nil {
		return nil, nil, err
	}

	var directChatroom *entity.ChatroomPopulated
	for _, cr := range chatrooms {
		if len(cr.Participants) == 2 {
			hasReceiver := false
			for _, p := range cr.Participants {
				if p.User.ID == params.ReceiverID {
					hasReceiver = true
					break
				}
			}
			if hasReceiver {
				directChatroom = cr
				break
			}
		}
	}

	var resultChatroom *entity.Chatroom

	// Create new direct chatroom if it doesn't exist
	if directChatroom == nil {
		createParams := chatroom.CreateChatroomParams{
			Type:         chatroomType,
			IsGroup:      false,
			Creator:      params.SenderID,
			Participants: []string{params.SenderID, params.ReceiverID},
		}

		var err error
		resultChatroom, err = s.chatroomService.CreateChatroom(ctx, createParams)
		if err != nil {
			return nil, nil, err
		}

		// Get the populated version of the new chatroom
		directChatroom, err = s.chatroomService.GetChatroom(ctx, resultChatroom.ID)
		if err != nil {
			return nil, nil, err
		}
	} else {
		// Convert populated chatroom to regular chatroom
		resultChatroom = &entity.Chatroom{
			ID:                   directChatroom.ID,
			Name:                 directChatroom.Name,
			IsGroup:              directChatroom.IsGroup,
			Type:                 directChatroom.Type,
			CreatedTimestamp:     directChatroom.CreatedTimestamp,
			LastMessage:          directChatroom.LastMessage,
			LastMessageTimestamp: directChatroom.LastMessageTimestamp,
			MessagesCount:        directChatroom.MessagesCount,
		}
	}

	// Send the message
	msgParams := SendMessageParams{
		ChatroomID: directChatroom.ID,
		SenderID:   params.SenderID,
		Message:    params.Message,
	}

	newChat, err := s.SendMessage(ctx, msgParams)
	if err != nil {
		return nil, nil, err
	}

	return newChat, resultChatroom, nil
}
