package repository

import (
	"app/pkg/quiz/domain/entity"
	"context"
)

// QuestionRepository defines the interface for question data access using GORM
type QuestionRepository interface {
	// FindOne retrieves a single question by ID
	FindOne(ctx context.Context, id uint) (*entity.Question, error)

	// FindAll retrieves multiple questions with filtering and pagination
	FindAll(ctx context.Context, query entity.QuestionQuery) ([]*entity.Question, int64, error)

	// Create stores a new question
	Create(ctx context.Context, question *entity.Question) error

	// Update modifies an existing question
	Update(ctx context.Context, question *entity.Question) error

	// Delete removes a question
	Delete(ctx context.Context, id uint) error
}
