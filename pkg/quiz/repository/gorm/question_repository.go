package gorm

import (
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/domain/repository"
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// QuestionRepository implements repository.QuestionRepository for GORM
type QuestionRepository struct {
	db *gorm.DB
}

// NewQuestionRepository creates a new GORM question repository
func NewQuestionRepository(db *gorm.DB) repository.QuestionRepository {
	return &QuestionRepository{
		db: db,
	}
}

// FindOne retrieves a single question by ID
func (r *QuestionRepository) FindOne(ctx context.Context, id uint) (*entity.Question, error) {
	var question entity.Question
	tx := r.db.WithContext(ctx).Preload("Answers").First(&question, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, tx.Error
	}

	return &question, nil
}

// FindAll retrieves multiple questions with filtering and pagination
func (r *QuestionRepository) FindAll(ctx context.Context, query entity.QuestionQuery) ([]*entity.Question, int64, error) {
	var questions []*entity.Question
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.Question{})

	// Apply filters
	if query.QuizID != 0 {
		db = db.Where("quiz_id = ?", query.QuizID)
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
	err := db.Order("created_at DESC").Preload("Answers").Find(&questions).Error
	if err != nil {
		return nil, 0, err
	}

	return questions, total, nil
}

// Create stores a new question
func (r *QuestionRepository) Create(ctx context.Context, question *entity.Question) error {
	return r.db.WithContext(ctx).Create(question).Error
}

// Update modifies an existing question
func (r *QuestionRepository) Update(ctx context.Context, question *entity.Question) error {
	return r.db.WithContext(ctx).Save(question).Error
}

// Delete removes a question
func (r *QuestionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Question{}, id).Error
}
