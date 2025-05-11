package entity

import (
	"app/pkg/types/pagination"
	"time"
)

// Question represents a quiz question entity
type Question struct {
	ID     uint   `json:"id" gorm:"primarykey, autoIncrement"`
	Text   string `json:"text"`
	QuizID uint   `json:"quizId" gorm:"column:quiz_id"`

	Answers []Answer `json:"answers" gorm:"constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// QuestionDTO represents the data transfer object for creating or updating a question
type QuestionDTO struct {
	Text    string      `json:"text" validate:"required"`
	QuizID  uint        `json:"quizId" validate:"required"`
	Answers []AnswerDTO `json:"answers" validate:"omitempty,dive"`
}

// QuestionQuery represents the query parameters for filtering questions
type QuestionQuery struct {
	pagination.Pagination
	QuizID  uint   `query:"quizId"`
	Keyword string `query:"keyword"`
}

// GetLimit returns the pagination limit or default value
func (q QuestionQuery) GetLimit() int {
	if q.Limit <= 0 {
		return 10 // Default limit
	}
	return q.Limit
}

// GetOffset returns the pagination offset
func (q QuestionQuery) GetOffset() int {
	return (q.Page - 1) * q.GetLimit()
}
