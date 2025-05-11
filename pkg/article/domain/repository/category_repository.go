package repository

import (
	"app/pkg/article/domain/entity"
	"context"
)

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	// FindOne retrieves a single category by ID
	FindOne(ctx context.Context, id uint) (*entity.Category, error)

	// FindBySlug retrieves a single category by its slug
	FindBySlug(ctx context.Context, slug string) (*entity.Category, error)

	// FindAll retrieves multiple categories with filtering and pagination
	FindAll(ctx context.Context, query entity.CategoryQuery) ([]*entity.Category, int64, error)

	// FindChildren retrieves all direct child categories of a parent category
	FindChildren(ctx context.Context, parentID uint) ([]*entity.Category, error)

	// GetHierarchy retrieves the complete ancestor hierarchy of a category
	GetHierarchy(ctx context.Context, id uint) ([]*entity.Category, error)

	// Create stores a new category
	Create(ctx context.Context, category *entity.Category) error

	// Update modifies an existing category
	Update(ctx context.Context, category *entity.Category) error

	// Delete removes a category
	Delete(ctx context.Context, id uint) error

	// CountBySlug counts categories with the given slug, excluding the specified ID if provided
	CountBySlug(ctx context.Context, slug string, excludeID *uint) (int64, error)
}
