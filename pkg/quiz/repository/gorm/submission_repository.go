package gorm

import (
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/domain/repository"
	"context"

	"gorm.io/gorm"
)

// SubmissionRepository implements repository.SubmissionRepository for GORM
type SubmissionRepository struct {
	db *gorm.DB
}

// NewSubmissionRepository creates a new GORM submission repository
func NewSubmissionRepository(db *gorm.DB) repository.SubmissionRepository {
	return &SubmissionRepository{
		db: db,
	}
}

// FindOne retrieves a single submission by ID
func (r *SubmissionRepository) FindOne(ctx context.Context, id uint) (*entity.Submission, error) {
	var submission entity.Submission
	tx := r.db.WithContext(ctx).Preload("Quiz").Preload("Question").Preload("Answer").Preload("User").First(&submission, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, tx.Error
	}

	return &submission, nil
}

// FindAll retrieves multiple submissions with filtering and pagination
func (r *SubmissionRepository) FindAll(ctx context.Context, query entity.SubmissionQuery) ([]*entity.Submission, int64, error) {
	var submissions []*entity.Submission
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.Submission{})

	// Apply filters
	if query.Quiz != 0 {
		db = db.Where("quiz_id = ?", query.Quiz)
	}

	if query.Question != 0 {
		db = db.Where("question_id = ?", query.Question)
	}

	if query.Answer != 0 {
		db = db.Where("answer_id = ?", query.Answer)
	}

	if query.User != 0 {
		db = db.Where("user_id = ?", query.User)
	}

	// Count total filtered records
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if query.Page > 0 && query.Limit > 0 {
		offset := (query.Page - 1) * query.Limit
		db = db.Offset(offset).Limit(query.Limit)
	}

	// Execute query with sorting and preloading related entities
	err := db.Order("created_at DESC").
		Preload("Quiz").
		Preload("Question").
		Preload("Answer").
		Preload("User").
		Find(&submissions).Error
	if err != nil {
		return nil, 0, err
	}

	return submissions, total, nil
}

// Create stores a new submission
func (r *SubmissionRepository) Create(ctx context.Context, submission *entity.Submission) error {
	return r.db.WithContext(ctx).Create(submission).Error
}

// CreateBulk stores multiple submissions in a transaction
func (r *SubmissionRepository) CreateBulk(ctx context.Context, submissions []*entity.Submission) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, submission := range submissions {
			if err := tx.Create(submission).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Update modifies an existing submission
func (r *SubmissionRepository) Update(ctx context.Context, submission *entity.Submission) error {
	return r.db.WithContext(ctx).Save(submission).Error
}

// Delete removes a submission
func (r *SubmissionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Submission{}, id).Error
}
