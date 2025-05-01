package repository

import (
	"app/pkg/quiz/domain/entity"
	"context"
)

// QuizRepository defines the interface for quiz data access using GORM
type QuizRepository interface {
	// FindOne retrieves a single quiz by ID
	FindOne(ctx context.Context, id uint) (*entity.Quiz, error)

	// FindAll retrieves multiple quizzes with filtering and pagination
	FindAll(ctx context.Context, query entity.QuizQuery) ([]*entity.Quiz, int64, error)

	// Create stores a new quiz
	Create(ctx context.Context, quiz *entity.Quiz) error

	// Update modifies an existing quiz
	Update(ctx context.Context, quiz *entity.Quiz) error

	// Delete removes a quiz
	Delete(ctx context.Context, id uint) error
}
