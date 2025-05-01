package answer

import (
	"app/pkg/quiz/domain/entity"
	"app/pkg/types/pagination"
	"context"
)

// AnswerService defines the interface for answer-related operations
type AnswerService interface {
	// FindOne retrieves an answer by ID
	FindOne(ctx context.Context, id uint) (*entity.Answer, error)

	// FindAll retrieves multiple answers with filtering and pagination
	FindAll(ctx context.Context, query entity.AnswerQuery) (*pagination.PaginatedResult[entity.Answer], error)

	// Create creates a new answer
	Create(ctx context.Context, answerDTO entity.AnswerDTO) (*entity.Answer, error)

	// Update modifies an existing answer
	Update(ctx context.Context, id uint, answerDTO entity.AnswerDTO) (*entity.Answer, error)

	// Delete removes an answer
	Delete(ctx context.Context, id uint) error
}
