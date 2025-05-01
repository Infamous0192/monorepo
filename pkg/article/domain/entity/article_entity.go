package entity

import (
	"app/pkg/types/pagination"
	"time"
)

// Article represents an article entity
type Article struct {
	ID          uint       `json:"id" gorm:"primarykey,autoIncrement"`
	Title       string     `json:"title"`
	Content     string     `json:"content"` // HTML content
	Slug        string     `json:"slug" gorm:"uniqueIndex"`
	PublishedAt *time.Time `json:"publishedAt,omitempty" gorm:"column:published_at"`

	Categories []*Category `json:"categories,omitempty" gorm:"many2many:article_categories;"`
	Tags       []*Tag      `json:"tags,omitempty" gorm:"many2many:article_tags;"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// ArticleDTO represents the data transfer object for creating or updating an article
type ArticleDTO struct {
	Title       string     `json:"title" validate:"required"`
	Content     string     `json:"content" validate:"required"`
	Slug        string     `json:"slug,omitempty"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
	CategoryIDs []uint     `json:"categoryIds,omitempty"`
	TagIDs      []uint     `json:"tagIds,omitempty"`
}

// ArticleQuery represents the query parameters for filtering articles
type ArticleQuery struct {
	pagination.Pagination
	Keyword     string `query:"keyword"`
	CategoryID  uint   `query:"categoryId"`
	TagID       uint   `query:"tagId"`
	IsPublished *bool  `query:"isPublished"`
}

// GetLimit returns the pagination limit or default value
func (q ArticleQuery) GetLimit() int {
	if q.Limit <= 0 {
		return 10 // Default limit
	}
	return q.Limit
}

// GetOffset returns the pagination offset
func (q ArticleQuery) GetOffset() int {
	return (q.Page - 1) * q.GetLimit()
}
