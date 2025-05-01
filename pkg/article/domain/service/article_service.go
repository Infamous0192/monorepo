package service

import (
	"app/pkg/article/domain/entity"
	"context"
)

// ArticleService defines the interface for article business logic
type ArticleService interface {
	// Get retrieves a single article by ID
	Get(ctx context.Context, id uint) (*entity.Article, error)

	// GetBySlug retrieves a single article by its slug
	GetBySlug(ctx context.Context, slug string) (*entity.Article, error)

	// GetAll retrieves multiple articles with filtering and pagination
	GetAll(ctx context.Context, query entity.ArticleQuery) ([]*entity.Article, int64, error)

	// Create stores a new article
	Create(ctx context.Context, dto entity.ArticleDTO) (*entity.Article, error)

	// Update modifies an existing article
	Update(ctx context.Context, id uint, dto entity.ArticleDTO) (*entity.Article, error)

	// Delete removes an article
	Delete(ctx context.Context, id uint) error

	// Publish sets an article as published
	Publish(ctx context.Context, id uint) error

	// Unpublish removes the published status of an article
	Unpublish(ctx context.Context, id uint) error
}
