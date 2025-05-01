package entity

import (
	"app/pkg/types/pagination"
	"time"
)

// Category represents a content category entity with hierarchical structure
type Category struct {
	ID          uint   `json:"id" gorm:"primarykey,autoIncrement"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug" gorm:"uniqueIndex"`

	// Self-referencing relationship for nested categories
	ParentID *uint       `json:"parentId,omitempty" gorm:"column:parent_id;index"`
	Parent   *Category   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []*Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`

	// Many-to-many relationship with articles
	Articles []*Article `json:"articles,omitempty" gorm:"many2many:article_categories;"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// CategoryDTO represents the data transfer object for creating or updating a category
type CategoryDTO struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Slug        string `json:"slug,omitempty"`
	ParentID    *uint  `json:"parentId,omitempty"`
}

// CategoryQuery represents the query parameters for filtering categories
type CategoryQuery struct {
	pagination.Pagination
	Keyword  string `query:"keyword"`
	ParentID *uint  `query:"parentId"`
}

// GetLimit returns the pagination limit or default value
func (q CategoryQuery) GetLimit() int {
	if q.Limit <= 0 {
		return 10 // Default limit
	}
	return q.Limit
}

// GetOffset returns the pagination offset
func (q CategoryQuery) GetOffset() int {
	return (q.Page - 1) * q.GetLimit()
}
