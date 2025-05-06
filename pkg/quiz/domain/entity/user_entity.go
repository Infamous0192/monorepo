package entity

import (
	"app/pkg/types/pagination"
	"time"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// User represents a user entity
type User struct {
	ID       uint   `json:"id" gorm:"primarykey, autoIncrement"`
	Name     string `json:"name" gorm:"type:varchar(100)"`
	Username string `json:"username" gorm:"type:varchar(100)"`
	Password string `json:"-" gorm:"type:varchar(200)"`
	Role     string `json:"role" gorm:"type:varchar(20);check:role IN ('admin', 'user')"`
	Status   bool   `json:"status"`

	BirthDate time.Time `json:"birthDate"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// UserDTO represents the data transfer object for creating or updating a user
type UserDTO struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required"`
	Status   bool   `json:"status" validate:"required"`

	BirthDate time.Time `json:"birthDate" validate:"required"`
}

// UserQuery represents the query parameters for filtering users
type UserQuery struct {
	pagination.Pagination
	Keyword string `query:"keyword"`
}

// GetLimit returns the pagination limit
func (q UserQuery) GetLimit() int {
	if q.Limit <= 0 {
		return 10 // Default limit
	}
	return q.Limit
}

// GetOffset returns the pagination offset
func (q UserQuery) GetOffset() int {
	return (q.Page - 1) * q.GetLimit()
}
