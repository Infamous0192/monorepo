package service

import (
	"app/pkg/article/domain/entity"
	"context"
)

// CategoryService defines the interface for category business logic
type CategoryService interface {
	// Get retrieves a single category by ID
	Get(ctx context.Context, id uint) (*entity.Category, error)

	// GetBySlug retrieves a single category by its slug
	GetBySlug(ctx context.Context, slug string) (*entity.Category, error)

	// GetAll retrieves multiple categories with filtering and pagination
	GetAll(ctx context.Context, query entity.CategoryQuery) ([]*entity.Category, int64, error)

	// GetHierarchy retrieves the complete hierarchy of a category
	GetHierarchy(ctx context.Context, id uint) ([]*entity.Category, error)

	// GetChildren retrieves all direct child categories of a parent
	GetChildren(ctx context.Context, parentID uint) ([]*entity.Category, error)

	// Create stores a new category
	Create(ctx context.Context, dto entity.CategoryDTO) (*entity.Category, error)

	// Update modifies an existing category
	Update(ctx context.Context, id uint, dto entity.CategoryDTO) (*entity.Category, error)

	// Delete removes a category
	Delete(ctx context.Context, id uint) error
}
