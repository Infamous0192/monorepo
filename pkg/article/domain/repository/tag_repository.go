package repository

import (
	"app/pkg/article/domain/entity"
	"context"
)

// TagRepository defines the interface for tag data access
type TagRepository interface {
	// FindOne retrieves a single tag by ID
	FindOne(ctx context.Context, id uint) (*entity.Tag, error)

	// FindBySlug retrieves a single tag by its slug
	FindBySlug(ctx context.Context, slug string) (*entity.Tag, error)

	// FindAll retrieves multiple tags with filtering and pagination
	FindAll(ctx context.Context, query entity.TagQuery) ([]*entity.Tag, int64, error)

	// FindByArticleID retrieves all tags associated with an article
	FindByArticleID(ctx context.Context, articleID uint) ([]*entity.Tag, error)

	// Create stores a new tag
	Create(ctx context.Context, tag *entity.Tag) error

	// Update modifies an existing tag
	Update(ctx context.Context, tag *entity.Tag) error

	// Delete removes a tag
	Delete(ctx context.Context, id uint) error
}
