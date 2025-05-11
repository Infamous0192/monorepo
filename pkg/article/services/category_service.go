package services

import (
	"app/pkg/article/domain/entity"
	"app/pkg/article/domain/repository"
	"app/pkg/article/domain/service"
	"app/pkg/exception"
	"context"
	"fmt"
	"strings"

	"github.com/gosimple/slug"
)

// categoryService implements the CategoryService interface
type categoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo repository.CategoryRepository) service.CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

// Get retrieves a single category by ID
func (s *categoryService) Get(ctx context.Context, id uint) (*entity.Category, error) {
	return s.categoryRepo.FindOne(ctx, id)
}

// GetBySlug retrieves a single category by its slug
func (s *categoryService) GetBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	return s.categoryRepo.FindBySlug(ctx, slug)
}

// GetAll retrieves multiple categories with filtering and pagination
func (s *categoryService) GetAll(ctx context.Context, query entity.CategoryQuery) ([]*entity.Category, int64, error) {
	return s.categoryRepo.FindAll(ctx, query)
}

// GetHierarchy retrieves the complete hierarchy of a category
func (s *categoryService) GetHierarchy(ctx context.Context, id uint) ([]*entity.Category, error) {
	return s.categoryRepo.GetHierarchy(ctx, id)
}

// GetChildren retrieves all direct child categories of a parent
func (s *categoryService) GetChildren(ctx context.Context, parentID uint) ([]*entity.Category, error) {
	return s.categoryRepo.FindChildren(ctx, parentID)
}

// generateUniqueSlug creates a unique slug by appending a numeric suffix if needed
func (s *categoryService) generateUniqueSlug(ctx context.Context, baseSlug string, excludeID *uint) (string, error) {
	// Try the base slug first
	uniqueSlug := baseSlug
	suffix := 1

	for {
		// Check if slug exists using count
		count, err := s.categoryRepo.CountBySlug(ctx, uniqueSlug, excludeID)
		if err != nil {
			return "", err
		}

		// If count is 0, slug is unique
		if count == 0 {
			return uniqueSlug, nil
		}

		// Try next suffix
		uniqueSlug = fmt.Sprintf("%s-%d", baseSlug, suffix)
		suffix++
	}
}

// Create stores a new category
func (s *categoryService) Create(ctx context.Context, dto entity.CategoryDTO) (*entity.Category, error) {
	// Validate name
	if strings.TrimSpace(dto.Name) == "" {
		return nil, exception.InvalidPayload(map[string]string{
			"name": "Category name is required",
		})
	}

	// Generate slug if not provided
	categorySlug := dto.Slug
	if categorySlug == "" {
		categorySlug = slug.Make(dto.Name)
	}

	// Ensure slug is unique
	uniqueSlug, err := s.generateUniqueSlug(ctx, categorySlug, nil)
	if err != nil {
		return nil, err
	}

	// Create category entity
	category := &entity.Category{
		Name:        dto.Name,
		Description: dto.Description,
		Slug:        uniqueSlug,
		ParentID:    dto.ParentID,
	}

	// Store category
	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	// Return created category with all relationships loaded
	return s.categoryRepo.FindOne(ctx, category.ID)
}

// Update modifies an existing category
func (s *categoryService) Update(ctx context.Context, id uint, dto entity.CategoryDTO) (*entity.Category, error) {
	// Check category exists
	category, err := s.categoryRepo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate name
	if strings.TrimSpace(dto.Name) == "" {
		return nil, exception.InvalidPayload(map[string]string{
			"name": "Category name is required",
		})
	}

	// Update fields
	category.Name = dto.Name
	category.Description = dto.Description

	// Generate slug if not provided
	if dto.Slug != "" {
		category.Slug = dto.Slug
	} else if category.Name != dto.Name {
		// Generate new slug if name changed and slug not specified
		baseSlug := slug.Make(dto.Name)
		uniqueSlug, err := s.generateUniqueSlug(ctx, baseSlug, &id)
		if err != nil {
			return nil, err
		}
		category.Slug = uniqueSlug
	}

	// Update parent if provided
	category.ParentID = dto.ParentID

	// Update category
	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	// Return updated category with all relationships loaded
	return s.categoryRepo.FindOne(ctx, category.ID)
}

// Delete removes a category
func (s *categoryService) Delete(ctx context.Context, id uint) error {
	return s.categoryRepo.Delete(ctx, id)
}
