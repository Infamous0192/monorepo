package repository

import (
	"app/pkg/quiz/domain/entity"
	"context"
)

// SubmissionRepository defines the interface for submission data access using GORM
type SubmissionRepository interface {
	// FindOne retrieves a single submission by ID
	FindOne(ctx context.Context, id uint) (*entity.Submission, error)

	// FindAll retrieves multiple submissions with filtering and pagination
	FindAll(ctx context.Context, query entity.SubmissionQuery) ([]*entity.Submission, int64, error)

	// Create stores a new submission
	Create(ctx context.Context, submission *entity.Submission) error

	// CreateBulk stores multiple submissions in a transaction
	CreateBulk(ctx context.Context, submissions []*entity.Submission) error

	// Update modifies an existing submission
	Update(ctx context.Context, submission *entity.Submission) error

	// Delete removes a submission
	Delete(ctx context.Context, id uint) error
}
