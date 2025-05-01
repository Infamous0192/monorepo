package user

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/chat/domain/repository"
	"app/pkg/types/pagination"
	"context"
)

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// Get retrieves a user by ID
func (s *userService) Get(ctx context.Context, id string) (*entity.User, error) {
	return s.userRepo.Get(ctx, id)
}

// GetAll retrieves multiple users with pagination
func (s *userService) GetAll(ctx context.Context, pag pagination.Pagination) ([]*entity.User, int64, error) {
	return s.userRepo.GetAll(ctx, pag)
}

// Create creates a new user
func (s *userService) Create(ctx context.Context, user *entity.User) error {
	return s.userRepo.Create(ctx, user)
}

// Update modifies an existing user
func (s *userService) Update(ctx context.Context, user *entity.User) error {
	return s.userRepo.Update(ctx, user)
}

// Delete removes a user
func (s *userService) Delete(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}
