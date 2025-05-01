package repository

import (
	"app/pkg/quiz/domain/entity"
	"context"
)

// AnswerRepository defines the interface for answer data access using GORM
type AnswerRepository interface {
	// FindOne retrieves a single answer by ID
	FindOne(ctx context.Context, id uint) (*entity.Answer, error)

	// FindAll retrieves multiple answers with filtering and pagination
	FindAll(ctx context.Context, query entity.AnswerQuery) ([]*entity.Answer, int64, error)

	// Create stores a new answer
	Create(ctx context.Context, answer *entity.Answer) error

	// Update modifies an existing answer
	Update(ctx context.Context, answer *entity.Answer) error

	// Delete removes an answer
	Delete(ctx context.Context, id uint) error
}
