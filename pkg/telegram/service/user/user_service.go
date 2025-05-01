package user

import (
	"app/pkg/exception"
	"app/pkg/telegram/domain/entity"
	"app/pkg/telegram/domain/repository"
	"app/pkg/types/pagination"
	"context"
)

type userService struct {
	userRepo   repository.UserRepository
	clientRepo repository.ClientRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository, clientRepo repository.ClientRepository) UserService {
	return &userService{
		userRepo:   userRepo,
		clientRepo: clientRepo,
	}
}

// Get retrieves a user by ID
func (s *userService) Get(ctx context.Context, id string) (*entity.User, error) {
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByTelegramID retrieves a user by Telegram ID
func (s *userService) GetByTelegramID(ctx context.Context, telegramID int64) (*entity.User, error) {
	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByClientID retrieves users by client ID with pagination
func (s *userService) GetByClientID(ctx context.Context, clientID string, pag pagination.Pagination) ([]*entity.User, int64, error) {
	// Verify client exists
	client, err := s.clientRepo.Get(ctx, clientID)
	if err != nil {
		return nil, 0, err
	}
	if client == nil {
		return nil, 0, exception.NotFound("Bot client")
	}

	return s.userRepo.GetByClientID(ctx, clientID, pag)
}

// GetAll retrieves multiple users with pagination
func (s *userService) GetAll(ctx context.Context, pag pagination.Pagination) ([]*entity.User, int64, error) {
	return s.userRepo.GetAll(ctx, pag)
}

// Create stores a new user
func (s *userService) Create(ctx context.Context, user *entity.User) error {
	// Verify client exists
	client, err := s.clientRepo.Get(ctx, user.ClientID)
	if err != nil {
		return err
	}
	if client == nil {
		return exception.NotFound("Bot client")
	}

	// Check if user with same Telegram ID exists
	existing, err := s.userRepo.GetByTelegramID(ctx, user.TelegramID)
	if err != nil {
		return err
	}
	if existing != nil {
		return exception.Http(400, "User with this Telegram ID already exists")
	}

	if user.Status == "" {
		user.Status = "active"
	}

	return s.userRepo.Create(ctx, user)
}

// Update modifies an existing user
func (s *userService) Update(ctx context.Context, user *entity.User) error {
	// Verify user exists
	existing, err := s.userRepo.Get(ctx, user.ID)
	if err != nil {
		return err
	}

	if existing == nil {
		return exception.NotFound("User")
	}

	// If client ID changed, verify new client exists
	if existing.ClientID != user.ClientID {
		client, err := s.clientRepo.Get(ctx, user.ClientID)
		if err != nil {
			return err
		}
		if client == nil {
			return exception.NotFound("Bot client")
		}
	}

	return s.userRepo.Update(ctx, user)
}

// Delete removes a user
func (s *userService) Delete(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}
