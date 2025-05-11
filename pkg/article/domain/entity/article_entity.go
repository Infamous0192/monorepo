package entity

import (
	"app/pkg/types/pagination"
	"mime/multipart"
	"time"
)

// Article represents an article entity
type Article struct {
	ID          uint       `json:"id" gorm:"primarykey,autoIncrement"`
	Title       string     `json:"title"`
	Content     string     `json:"content"` // HTML content
	Slug        string     `json:"slug" gorm:"uniqueIndex"`
	PublishedAt *time.Time `json:"publishedAt,omitempty" gorm:"column:published_at"`
	ThumbnailID *uint      `json:"thumbnailId,omitempty"`
	Thumbnail   *File      `json:"thumbnail" gorm:"foreignKey:ThumbnailID"`

	Categories []*Category `json:"categories" gorm:"many2many:article_categories;"`
	Tags       []*Tag      `json:"tags" gorm:"many2many:article_tags;"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
}

// ArticleDTO represents the data transfer object for creating or updating an article
type ArticleDTO struct {
	Title       string     `form:"title" validate:"required"`
	Content     string     `form:"content" validate:"required"`
	Slug        string     `form:"slug,omitempty"`
	PublishedAt *time.Time `form:"publishedAt,omitempty"`
	CategoryIDs []uint     `form:"categoryIds,omitempty"`
	TagIDs      []uint     `form:"tagIds,omitempty"`
	ThumbnailID *uint      `form:"thumbnailId,omitempty"`

	Thumbnail *multipart.FileHeader `form:"thumbnail,omitempty"`
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
