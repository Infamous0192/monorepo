package entity

import (
	"app/pkg/types/pagination"
	"time"
)

// Submission represents a quiz submission entity
type Submission struct {
	ID uint `json:"id" gorm:"primarykey, autoIncrement"`

	Quiz   *Quiz `json:"quiz,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	QuizID uint  `json:"-"`

	Question   *Question `json:"question,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	QuestionID uint      `json:"-"`

	Answer   *Answer `json:"answer,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	AnswerID uint    `json:"-"`

	UserID uint  `json:"-"`
	User   *User `json:"user,omitempty" gorm:"constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// SubmissionDTO represents the data transfer object for creating a submission
type SubmissionDTO struct {
	Quiz     uint `json:"quiz" form:"quiz" validate:"required"`
	Question uint `json:"question" form:"question" validate:"required"`
	Answer   uint `json:"answer" form:"answer" validate:"required"`
	User     uint `json:"user" form:"user" validate:"required"`
}

// SubmissionInsertDTO represents a bulk submission DTO
type SubmissionInsertDTO struct {
	Data []SubmissionDTO `json:"data" form:"data" validate:"required,dive,required"`
}

// SubmissionQuery represents the query parameters for filtering submissions
type SubmissionQuery struct {
	pagination.Pagination
	Quiz     uint `query:"quiz"`
	Question uint `query:"question"`
	Answer   uint `query:"answer"`
	User     uint `query:"user"`
}

// GetLimit returns the pagination limit or default value
func (q SubmissionQuery) GetLimit() int {
	if q.Limit <= 0 {
		return 10 // Default limit
	}
	return q.Limit
}

// GetOffset returns the pagination offset
func (q SubmissionQuery) GetOffset() int {
	return (q.Page - 1) * q.GetLimit()
}
