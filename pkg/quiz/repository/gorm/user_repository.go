package gorm

import (
	"app/pkg/exception"
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/domain/repository"
	"errors"

	"gorm.io/gorm"
)

// userRepository is a GORM implementation of the user repository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new GORM user repository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

// Get retrieves a user by ID
func (r *userRepository) Get(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NotFound("user")
		}
		return nil, exception.InternalError("Failed to get user: " + err.Error())
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil without error when user not found
		}
		return nil, exception.InternalError("Failed to get user by username: " + err.Error())
	}
	return &user, nil
}

// GetAll retrieves multiple users with pagination and filtering
func (r *userRepository) GetAll(query entity.UserQuery) ([]*entity.User, int64, error) {
	var users []*entity.User
	var count int64

	db := r.db
	if query.Keyword != "" {
		db = db.Where("name LIKE ? OR username LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	// Get total count
	if err := db.Model(&entity.User{}).Count(&count).Error; err != nil {
		return nil, 0, exception.InternalError("Failed to count users: " + err.Error())
	}

	// Get paginated data
	if err := db.Limit(query.GetLimit()).Offset(query.GetOffset()).Find(&users).Error; err != nil {
		return nil, 0, exception.InternalError("Failed to get users: " + err.Error())
	}

	return users, count, nil
}

// Create stores a new user
func (r *userRepository) Create(user *entity.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return exception.InternalError("Failed to create user: " + err.Error())
	}
	return nil
}

// Update modifies an existing user
func (r *userRepository) Update(user *entity.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return exception.InternalError("Failed to update user: " + err.Error())
	}
	return nil
}

// Delete removes a user
func (r *userRepository) Delete(id uint) error {
	if err := r.db.Delete(&entity.User{}, id).Error; err != nil {
		return exception.InternalError("Failed to delete user: " + err.Error())
	}
	return nil
}
