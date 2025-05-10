package entity

import (
	"time"
)

// File represents a file entity for storing uploaded files
type File struct {
	ID          uint      `json:"id" gorm:"primarykey,autoIncrement"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	ContentType string    `json:"contentType"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// FileDTO represents the data transfer object for file information
type FileDTO struct {
	Name        string `json:"name" validate:"required"`
	Path        string `json:"path" validate:"required"`
	Size        int64  `json:"size" validate:"required"`
	ContentType string `json:"contentType" validate:"required"`
}
