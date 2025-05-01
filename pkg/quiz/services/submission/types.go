package submission

import (
	"app/pkg/quiz/domain/entity"
	"app/pkg/types/pagination"
	"context"
)

// SubmissionService defines the interface for submission-related operations
type SubmissionService interface {
	// FindOne retrieves a submission by ID
	FindOne(ctx context.Context, id uint) (*entity.Submission, error)

	// FindAll retrieves multiple submissions with filtering and pagination
	FindAll(ctx context.Context, query entity.SubmissionQuery) (*pagination.PaginatedResult[entity.Submission], error)

	// Create creates a new submission
	Create(ctx context.Context, submissionDTO entity.SubmissionDTO) (*entity.Submission, error)

	// CreateBulk creates multiple submissions
	CreateBulk(ctx context.Context, submissionDTO entity.SubmissionInsertDTO) ([]*entity.Submission, error)

	// Update modifies an existing submission
	Update(ctx context.Context, id uint, submissionDTO entity.SubmissionDTO) (*entity.Submission, error)

	// Delete removes a submission
	Delete(ctx context.Context, id uint) error
}
