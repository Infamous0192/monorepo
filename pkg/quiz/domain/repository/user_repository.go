package repository

import (
	"app/pkg/quiz/domain/entity"
)

// UserRepository defines methods for user data persistence
type UserRepository interface {
	// Get retrieves a user by ID
	Get(id uint) (*entity.User, error)

	// GetByUsername retrieves a user by their username
	GetByUsername(username string) (*entity.User, error)

	// GetAll retrieves multiple users with pagination and filtering
	GetAll(query entity.UserQuery) ([]*entity.User, int64, error)

	// Create stores a new user
	Create(user *entity.User) error

	// Update modifies an existing user
	Update(user *entity.User) error

	// Delete removes a user
	Delete(id uint) error
}
