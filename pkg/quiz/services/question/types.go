package question

import (
	"app/pkg/quiz/domain/entity"
	"app/pkg/types/pagination"
	"context"
)

// QuestionService defines the interface for question-related operations
type QuestionService interface {
	// FindOne retrieves a question by ID
	FindOne(ctx context.Context, id uint) (*entity.Question, error)

	// FindAll retrieves multiple questions with filtering and pagination
	FindAll(ctx context.Context, query entity.QuestionQuery) (*pagination.PaginatedResult[entity.Question], error)

	// Create creates a new question
	Create(ctx context.Context, questionDTO entity.QuestionDTO) (*entity.Question, error)

	// Update modifies an existing question
	Update(ctx context.Context, id uint, questionDTO entity.QuestionDTO) (*entity.Question, error)

	// Delete removes a question
	Delete(ctx context.Context, id uint) error
}
