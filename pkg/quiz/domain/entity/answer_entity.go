package entity

import (
	"app/pkg/types/pagination"
	"time"
)

// Answer represents an answer entity
type Answer struct {
	ID    uint   `json:"id" gorm:"primarykey, autoIncrement"`
	Text  string `json:"text"`
	Value int32  `json:"value"`

	Question   *Question `json:"question,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	QuestionID uint      `json:"-"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// AnswerDTO represents the data transfer object for creating or updating an answer
type AnswerDTO struct {
	Text     string `json:"text" form:"text" validate:"required"`
	Value    *int32 `json:"value" form:"value" validate:"required"`
	Question uint   `json:"question" form:"question" validate:"omitempty,exist=questions"`
}

// AnswerQuery represents the query parameters for filtering answers
type AnswerQuery struct {
	pagination.Pagination
	Question uint   `query:"question"`
	Keyword  string `query:"keyword"`
}

// GetLimit returns the pagination limit or default value
func (q AnswerQuery) GetLimit() int {
	if q.Limit <= 0 {
		return 10 // Default limit
	}
	return q.Limit
}

// GetOffset returns the pagination offset
func (q AnswerQuery) GetOffset() int {
	return (q.Page - 1) * q.GetLimit()
}
