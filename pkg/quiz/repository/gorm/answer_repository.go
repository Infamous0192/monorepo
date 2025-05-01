package gorm

import (
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/domain/repository"
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// AnswerRepository implements repository.AnswerRepository for GORM
type AnswerRepository struct {
	db *gorm.DB
}

// NewAnswerRepository creates a new GORM answer repository
func NewAnswerRepository(db *gorm.DB) repository.AnswerRepository {
	return &AnswerRepository{
		db: db,
	}
}

// FindOne retrieves a single answer by ID
func (r *AnswerRepository) FindOne(ctx context.Context, id uint) (*entity.Answer, error) {
	var answer entity.Answer
	tx := r.db.WithContext(ctx).First(&answer, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, tx.Error
	}

	return &answer, nil
}

// FindAll retrieves multiple answers with filtering and pagination
func (r *AnswerRepository) FindAll(ctx context.Context, query entity.AnswerQuery) ([]*entity.Answer, int64, error) {
	var answers []*entity.Answer
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.Answer{})

	// Apply filters
	if query.Question != 0 {
		db = db.Where("question_id = ?", query.Question)
	}

	if query.Keyword != "" {
		searchTerm := fmt.Sprintf("%%%s%%", strings.ToLower(query.Keyword))
		db = db.Where("LOWER(text) LIKE ?", searchTerm)
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

	// Execute query with sorting
	err := db.Order("created_at DESC").Find(&answers).Error
	if err != nil {
		return nil, 0, err
	}

	return answers, total, nil
}

// Create stores a new answer
func (r *AnswerRepository) Create(ctx context.Context, answer *entity.Answer) error {
	return r.db.WithContext(ctx).Create(answer).Error
}

// Update modifies an existing answer
func (r *AnswerRepository) Update(ctx context.Context, answer *entity.Answer) error {
	return r.db.WithContext(ctx).Save(answer).Error
}

// Delete removes an answer
func (r *AnswerRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Answer{}, id).Error
}
