package user

import (
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/domain/repository"
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

// FindOne retrieves a user by ID
func (s *userService) FindOne(ctx context.Context, id uint) (*entity.User, error) {
	return s.userRepo.Get(id)
}

// FindAll retrieves multiple users with filtering and pagination
func (s *userService) FindAll(ctx context.Context, query entity.UserQuery) (*pagination.PaginatedResult[entity.User], error) {
	items, total, err := s.userRepo.GetAll(query)
	if err != nil {
		return nil, err
	}

	// Calculate pagination metadata
	hasNext := (int64(query.Page*query.GetLimit()) < total)
	hasPrev := query.Page > 1

	result := &pagination.PaginatedResult[entity.User]{
		Metadata: pagination.Metadata{
			Pagination: pagination.Pagination{
				Page:  query.Page,
				Limit: query.GetLimit(),
			},
			Total:   total,
			Count:   len(items),
			HasNext: hasNext,
			HasPrev: hasPrev,
		},
		Result: items,
	}

	return result, nil
}

// Create creates a new user
func (s *userService) Create(ctx context.Context, userDTO entity.UserDTO) (*entity.User, error) {
	// Create user from DTO
	user := &entity.User{
		Name:      userDTO.Name,
		Username:  userDTO.Username,
		Password:  userDTO.Password, // Note: This should be hashed before saving
		Role:      userDTO.Role,
		Status:    userDTO.Status,
		BirthDate: userDTO.BirthDate,
	}

	// Create user in repository
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Update modifies an existing user
func (s *userService) Update(ctx context.Context, id uint, userDTO entity.UserDTO) (*entity.User, error) {
	// Get existing user
	user, err := s.userRepo.Get(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	// Update user properties
	user.Name = userDTO.Name
	user.Username = userDTO.Username
	user.Role = userDTO.Role
	user.Status = userDTO.Status
	user.BirthDate = userDTO.BirthDate

	// Only update password if provided in DTO
	if userDTO.Password != "" {
		user.Password = userDTO.Password // Note: This should be hashed before saving
	}

	// Update user in repository
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Delete removes a user
func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.userRepo.Delete(id)
}
