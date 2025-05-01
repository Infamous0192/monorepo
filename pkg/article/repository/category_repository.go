package repository

import (
	"app/pkg/article/domain/entity"
	"app/pkg/article/domain/repository"
	"app/pkg/exception"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// categoryRepository implements the CategoryRepository interface using GORM
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new instance of the category repository
func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

// FindOne retrieves a single category by ID
func (r *categoryRepository) FindOne(ctx context.Context, id uint) (*entity.Category, error) {
	var category entity.Category

	result := r.db.WithContext(ctx).
		Preload("Children").
		First(&category, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, exception.NotFound("Category")
		}
		return nil, exception.InternalError(fmt.Sprintf("Failed to find category: %v", result.Error))
	}

	return &category, nil
}

// FindBySlug retrieves a single category by its slug
func (r *categoryRepository) FindBySlug(ctx context.Context, slug string) (*entity.Category, error) {
	var category entity.Category

	result := r.db.WithContext(ctx).
		Preload("Children").
		Where("slug = ?", slug).
		First(&category)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, exception.NotFound("Category")
		}
		return nil, exception.InternalError(fmt.Sprintf("Failed to find category by slug: %v", result.Error))
	}

	return &category, nil
}

// FindAll retrieves multiple categories with filtering and pagination
func (r *categoryRepository) FindAll(ctx context.Context, query entity.CategoryQuery) ([]*entity.Category, int64, error) {
	var categories []*entity.Category
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.Category{})

	// Apply filters
	if query.Keyword != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	if query.ParentID != nil {
		if *query.ParentID == 0 {
			// Root categories (no parent)
			db = db.Where("parent_id IS NULL")
		} else {
			// Child categories of specified parent
			db = db.Where("parent_id = ?", *query.ParentID)
		}
	}

	// Count total records
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, exception.InternalError(fmt.Sprintf("Failed to count categories: %v", err))
	}

	// Apply pagination and retrieve data
	result := db.
		Preload("Children").
		Limit(query.GetLimit()).
		Offset(query.GetOffset()).
		Order("name ASC").
		Find(&categories)

	if result.Error != nil {
		return nil, 0, exception.InternalError(fmt.Sprintf("Failed to find categories: %v", result.Error))
	}

	return categories, total, nil
}

// FindChildren retrieves all direct child categories of a parent category
func (r *categoryRepository) FindChildren(ctx context.Context, parentID uint) ([]*entity.Category, error) {
	var categories []*entity.Category

	result := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("name ASC").
		Find(&categories)

	if result.Error != nil {
		return nil, exception.InternalError(fmt.Sprintf("Failed to find child categories: %v", result.Error))
	}

	return categories, nil
}

// GetHierarchy retrieves the complete ancestor hierarchy of a category
func (r *categoryRepository) GetHierarchy(ctx context.Context, id uint) ([]*entity.Category, error) {
	var hierarchy []*entity.Category

	// First, check if the category exists
	category, err := r.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	// Add the current category to the hierarchy
	hierarchy = append(hierarchy, category)

	// If the category has a parent, recursively fetch its ancestors
	currentID := category.ParentID
	for currentID != nil && *currentID > 0 {
		parent, err := r.FindOne(ctx, *currentID)
		if err != nil {
			return nil, exception.InternalError(fmt.Sprintf("Failed to get parent category: %v", err))
		}

		// Add parent to the beginning of the hierarchy
		hierarchy = append([]*entity.Category{parent}, hierarchy...)

		// Move up to the next parent
		currentID = parent.ParentID
	}

	return hierarchy, nil
}

// Create stores a new category
func (r *categoryRepository) Create(ctx context.Context, category *entity.Category) error {
	// Check if slug already exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Category{}).
		Where("slug = ?", category.Slug).
		Count(&count).Error; err != nil {
		return exception.InternalError(fmt.Sprintf("Failed to check slug uniqueness: %v", err))
	}

	if count > 0 {
		return exception.InvalidPayload(map[string]string{
			"slug": "Slug already exists",
		})
	}

	// Check if parent exists if specified
	if category.ParentID != nil && *category.ParentID > 0 {
		var parentCount int64
		if err := r.db.WithContext(ctx).Model(&entity.Category{}).
			Where("id = ?", *category.ParentID).
			Count(&parentCount).Error; err != nil {
			return exception.InternalError(fmt.Sprintf("Failed to check parent category: %v", err))
		}

		if parentCount == 0 {
			return exception.InvalidPayload(map[string]string{
				"parentId": "Parent category does not exist",
			})
		}
	}

	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return exception.InternalError(fmt.Sprintf("Failed to create category: %v", err))
	}

	return nil
}

// Update modifies an existing category
func (r *categoryRepository) Update(ctx context.Context, category *entity.Category) error {
	// Check if category exists
	var existingCategory entity.Category
	if err := r.db.WithContext(ctx).First(&existingCategory, category.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return exception.NotFound("Category")
		}
		return exception.InternalError(fmt.Sprintf("Failed to find category: %v", err))
	}

	// Check slug uniqueness if changed
	if category.Slug != existingCategory.Slug {
		var count int64
		if err := r.db.WithContext(ctx).Model(&entity.Category{}).
			Where("slug = ? AND id != ?", category.Slug, category.ID).
			Count(&count).Error; err != nil {
			return exception.InternalError(fmt.Sprintf("Failed to check slug uniqueness: %v", err))
		}

		if count > 0 {
			return exception.InvalidPayload(map[string]string{
				"slug": "Slug already exists",
			})
		}
	}

	// Prevent category from becoming its own parent or child
	if category.ParentID != nil && *category.ParentID > 0 {
		// Check parent exists
		var parentCount int64
		if err := r.db.WithContext(ctx).Model(&entity.Category{}).
			Where("id = ?", *category.ParentID).
			Count(&parentCount).Error; err != nil {
			return exception.InternalError(fmt.Sprintf("Failed to check parent category: %v", err))
		}

		if parentCount == 0 {
			return exception.InvalidPayload(map[string]string{
				"parentId": "Parent category does not exist",
			})
		}

		// Check for circular reference
		if *category.ParentID == category.ID {
			return exception.InvalidPayload(map[string]string{
				"parentId": "A category cannot be its own parent",
			})
		}

		// Check if parent is one of the children (to avoid cycles)
		children, err := r.FindChildren(ctx, category.ID)
		if err != nil {
			return err
		}

		// Recursive function to check if potential parent is a descendant
		var isDescendant func(parentID uint, children []*entity.Category) bool
		isDescendant = func(parentID uint, children []*entity.Category) bool {
			for _, child := range children {
				if child.ID == parentID {
					return true
				}

				// If this child has children, check them too
				if len(child.Children) > 0 {
					if isDescendant(parentID, child.Children) {
						return true
					}
				} else {
					// If no loaded children, we need to fetch them
					childChildren, err := r.FindChildren(ctx, child.ID)
					if err == nil && len(childChildren) > 0 {
						if isDescendant(parentID, childChildren) {
							return true
						}
					}
				}
			}
			return false
		}

		if isDescendant(*category.ParentID, children) {
			return exception.InvalidPayload(map[string]string{
				"parentId": "Cannot set a child category as the parent (would create a cycle)",
			})
		}
	}

	// Update the category
	result := r.db.WithContext(ctx).Model(category).Updates(map[string]interface{}{
		"name":        category.Name,
		"description": category.Description,
		"slug":        category.Slug,
		"parent_id":   category.ParentID,
		"updated_at":  time.Now(),
	})

	if result.Error != nil {
		return exception.InternalError(fmt.Sprintf("Failed to update category: %v", result.Error))
	}

	return nil
}

// Delete removes a category
func (r *categoryRepository) Delete(ctx context.Context, id uint) error {
	// Check for children
	var childrenCount int64
	if err := r.db.WithContext(ctx).Model(&entity.Category{}).
		Where("parent_id = ?", id).
		Count(&childrenCount).Error; err != nil {
		return exception.InternalError(fmt.Sprintf("Failed to check child categories: %v", err))
	}

	if childrenCount > 0 {
		return exception.BadRequest("Cannot delete a category with child categories. Move or delete children first.")
	}

	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return exception.InternalError(fmt.Sprintf("Failed to begin transaction: %v", tx.Error))
	}

	// Remove category from articles' associations
	if err := tx.Exec("DELETE FROM article_categories WHERE category_id = ?", id).Error; err != nil {
		tx.Rollback()
		return exception.InternalError(fmt.Sprintf("Failed to remove category associations: %v", err))
	}

	// Delete the category
	result := tx.Delete(&entity.Category{}, id)
	if result.Error != nil {
		tx.Rollback()
		return exception.InternalError(fmt.Sprintf("Failed to delete category: %v", result.Error))
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return exception.NotFound("Category")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return exception.InternalError(fmt.Sprintf("Failed to commit transaction: %v", err))
	}

	return nil
}
