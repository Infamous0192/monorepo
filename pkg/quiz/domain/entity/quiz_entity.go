package entity

import (
	"app/pkg/types/pagination"
	"time"
)

// Quiz represents a quiz entity
type Quiz struct {
	ID          uint   `json:"id" gorm:"primarykey, autoIncrement"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Questions []Question `json:"questions"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// QuizDTO represents the data transfer object for creating or updating a quiz
type QuizDTO struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
}

// QuizQuery represents the query parameters for filtering quizzes
type QuizQuery struct {
	pagination.Pagination
	Keyword string `query:"keyword"`
}

// GetLimit returns the pagination limit or default value
func (q QuizQuery) GetLimit() int {
	if q.Limit <= 0 {
		return 10 // Default limit
	}
	return q.Limit
}

// GetOffset returns the pagination offset
func (q QuizQuery) GetOffset() int {
	return (q.Page - 1) * q.GetLimit()
}
