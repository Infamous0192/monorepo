package entity

import (
	"app/pkg/types/pagination"
	"time"
)

// Tag represents a tag entity for categorizing articles
type Tag struct {
	ID          uint   `json:"id" gorm:"primarykey,autoIncrement"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug" gorm:"uniqueIndex"`

	// Many-to-many relationship with articles
	Articles []*Article `json:"articles,omitempty" gorm:"many2many:article_tags;"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// TagDTO represents the data transfer object for creating or updating a tag
type TagDTO struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Slug        string `json:"slug,omitempty"`
}

// TagQuery represents the query parameters for filtering tags
type TagQuery struct {
	pagination.Pagination
	Keyword string `query:"keyword"`
}

// GetLimit returns the pagination limit or default value
func (q TagQuery) GetLimit() int {
	if q.Limit <= 0 {
		return 10 // Default limit
	}
	return q.Limit
}

// GetOffset returns the pagination offset
func (q TagQuery) GetOffset() int {
	return (q.Page - 1) * q.GetLimit()
}
