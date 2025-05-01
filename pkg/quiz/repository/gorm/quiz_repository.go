package gorm

import (
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/domain/repository"
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// QuizRepository implements repository.QuizRepository for GORM
type QuizRepository struct {
	db *gorm.DB
}

// NewQuizRepository creates a new GORM quiz repository
func NewQuizRepository(db *gorm.DB) repository.QuizRepository {
	return &QuizRepository{
		db: db,
	}
}

// FindOne retrieves a single quiz by ID
func (r *QuizRepository) FindOne(ctx context.Context, id uint) (*entity.Quiz, error) {
	var quiz entity.Quiz
	tx := r.db.WithContext(ctx).Preload("Questions").First(&quiz, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, tx.Error
	}

	return &quiz, nil
}

// FindAll retrieves multiple quizzes with filtering and pagination
func (r *QuizRepository) FindAll(ctx context.Context, query entity.QuizQuery) ([]*entity.Quiz, int64, error) {
	var quizzes []*entity.Quiz
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.Quiz{})

	// Apply filters
	if query.Keyword != "" {
		searchTerm := fmt.Sprintf("%%%s%%", strings.ToLower(query.Keyword))
		db = db.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)
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
	err := db.Order("created_at DESC").Preload("Questions").Find(&quizzes).Error
	if err != nil {
		return nil, 0, err
	}

	return quizzes, total, nil
}

// Create stores a new quiz
func (r *QuizRepository) Create(ctx context.Context, quiz *entity.Quiz) error {
	return r.db.WithContext(ctx).Create(quiz).Error
}

// Update modifies an existing quiz
func (r *QuizRepository) Update(ctx context.Context, quiz *entity.Quiz) error {
	return r.db.WithContext(ctx).Save(quiz).Error
}

// Delete removes a quiz
func (r *QuizRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Quiz{}, id).Error
}
