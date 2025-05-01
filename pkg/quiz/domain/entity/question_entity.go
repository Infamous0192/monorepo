package entity

import (
	"app/pkg/types/pagination"
	"time"
)

// Question represents a quiz question entity
type Question struct {
	ID      uint   `json:"id" gorm:"primarykey, autoIncrement"`
	Content string `json:"content"`
	QuizID  uint   `json:"quizId" gorm:"column:quiz_id"`

	Options []Option `json:"options"`
	Answers []Answer `json:"answers"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// Option represents a question option entity
type Option struct {
	ID         uint   `json:"id" gorm:"primarykey, autoIncrement"`
	Content    string `json:"content"`
	IsCorrect  bool   `json:"isCorrect" gorm:"column:is_correct"`
	QuestionID uint   `json:"questionId" gorm:"column:question_id"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// OptionDTO represents the data transfer object for creating or updating an option
type OptionDTO struct {
	Content   string `json:"content" validate:"required"`
	IsCorrect bool   `json:"isCorrect"`
}

// QuestionDTO represents the data transfer object for creating or updating a question
type QuestionDTO struct {
	Content string      `json:"content" validate:"required"`
	QuizID  uint        `json:"quizId" validate:"required"`
	Options []OptionDTO `json:"options" validate:"omitempty,dive"`
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
