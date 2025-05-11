package repository

import (
	"app/pkg/article/domain/entity"
	"context"
)

// ArticleRepository defines the interface for article data access
type ArticleRepository interface {
	// FindOne retrieves a single article by ID
	FindOne(ctx context.Context, id uint) (*entity.Article, error)

	// FindBySlug retrieves a single article by its slug
	FindBySlug(ctx context.Context, slug string) (*entity.Article, error)

	// FindAll retrieves multiple articles with filtering and pagination
	FindAll(ctx context.Context, query entity.ArticleQuery) ([]*entity.Article, int64, error)

	// Create stores a new article
	Create(ctx context.Context, article *entity.Article) error

	// Update modifies an existing article
	Update(ctx context.Context, article *entity.Article) error

	// Delete removes an article
	Delete(ctx context.Context, id uint) error

	// Publish sets the published status of an article
	Publish(ctx context.Context, id uint, publishedAt *string) error

	// CountBySlug counts articles with the given slug, excluding the specified ID if provided
	CountBySlug(ctx context.Context, slug string, excludeID *uint) (int64, error)
}
