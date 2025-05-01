package service

import (
	"app/pkg/article/domain/entity"
	"context"
)

// TagService defines the interface for tag business logic
type TagService interface {
	// Get retrieves a single tag by ID
	Get(ctx context.Context, id uint) (*entity.Tag, error)

	// GetBySlug retrieves a single tag by its slug
	GetBySlug(ctx context.Context, slug string) (*entity.Tag, error)

	// GetAll retrieves multiple tags with filtering and pagination
	GetAll(ctx context.Context, query entity.TagQuery) ([]*entity.Tag, int64, error)

	// GetByArticle retrieves all tags associated with an article
	GetByArticle(ctx context.Context, articleID uint) ([]*entity.Tag, error)

	// Create stores a new tag
	Create(ctx context.Context, dto entity.TagDTO) (*entity.Tag, error)

	// Update modifies an existing tag
	Update(ctx context.Context, id uint, dto entity.TagDTO) (*entity.Tag, error)

	// Delete removes a tag
	Delete(ctx context.Context, id uint) error
}
