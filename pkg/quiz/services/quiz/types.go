package quiz

import (
	"app/pkg/quiz/domain/entity"
	"app/pkg/types/pagination"
	"context"
)

// QuizService defines the interface for quiz-related operations
type QuizService interface {
	// FindOne retrieves a quiz by ID
	FindOne(ctx context.Context, id uint) (*entity.Quiz, error)

	// FindAll retrieves multiple quizzes with filtering and pagination
	FindAll(ctx context.Context, query entity.QuizQuery) (*pagination.PaginatedResult[entity.Quiz], error)

	// Create creates a new quiz
	Create(ctx context.Context, quizDTO entity.QuizDTO) (*entity.Quiz, error)

	// Update modifies an existing quiz
	Update(ctx context.Context, id uint, quizDTO entity.QuizDTO) (*entity.Quiz, error)

	// Delete removes a quiz
	Delete(ctx context.Context, id uint) error
}
