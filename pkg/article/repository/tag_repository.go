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

// tagRepository implements the TagRepository interface using GORM
type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository creates a new instance of the tag repository
func NewTagRepository(db *gorm.DB) repository.TagRepository {
	return &tagRepository{
		db: db,
	}
}

// FindOne retrieves a single tag by ID
func (r *tagRepository) FindOne(ctx context.Context, id uint) (*entity.Tag, error) {
	var tag entity.Tag

	result := r.db.WithContext(ctx).First(&tag, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, exception.NotFound("Tag")
		}
		return nil, exception.InternalError(fmt.Sprintf("Failed to find tag: %v", result.Error))
	}

	return &tag, nil
}

// FindBySlug retrieves a single tag by its slug
func (r *tagRepository) FindBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	var tag entity.Tag

	result := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		First(&tag)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, exception.NotFound("Tag")
		}
		return nil, exception.InternalError(fmt.Sprintf("Failed to find tag by slug: %v", result.Error))
	}

	return &tag, nil
}

// FindAll retrieves multiple tags with filtering and pagination
func (r *tagRepository) FindAll(ctx context.Context, query entity.TagQuery) ([]*entity.Tag, int64, error) {
	var tags []*entity.Tag
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.Tag{})

	// Apply filters
	if query.Keyword != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	// Count total records
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, exception.InternalError(fmt.Sprintf("Failed to count tags: %v", err))
	}

	// Apply pagination and retrieve data
	result := db.
		Limit(query.GetLimit()).
		Offset(query.GetOffset()).
		Order("name ASC").
		Find(&tags)

	if result.Error != nil {
		return nil, 0, exception.InternalError(fmt.Sprintf("Failed to find tags: %v", result.Error))
	}

	return tags, total, nil
}

// FindByArticleID retrieves all tags associated with an article
func (r *tagRepository) FindByArticleID(ctx context.Context, articleID uint) ([]*entity.Tag, error) {
	var tags []*entity.Tag

	result := r.db.WithContext(ctx).
		Joins("JOIN article_tags ON article_tags.tag_id = tags.id").
		Where("article_tags.article_id = ?", articleID).
		Order("name ASC").
		Find(&tags)

	if result.Error != nil {
		return nil, exception.InternalError(fmt.Sprintf("Failed to find article tags: %v", result.Error))
	}

	return tags, nil
}

// Create stores a new tag
func (r *tagRepository) Create(ctx context.Context, tag *entity.Tag) error {
	// Check if slug already exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Tag{}).
		Where("slug = ?", tag.Slug).
		Count(&count).Error; err != nil {
		return exception.InternalError(fmt.Sprintf("Failed to check slug uniqueness: %v", err))
	}

	if count > 0 {
		return exception.InvalidPayload(map[string]string{
			"slug": "Slug already exists",
		})
	}

	// Check if name already exists (tags should have unique names)
	if err := r.db.WithContext(ctx).Model(&entity.Tag{}).
		Where("name = ?", tag.Name).
		Count(&count).Error; err != nil {
		return exception.InternalError(fmt.Sprintf("Failed to check name uniqueness: %v", err))
	}

	if count > 0 {
		return exception.InvalidPayload(map[string]string{
			"name": "Tag name already exists",
		})
	}

	if err := r.db.WithContext(ctx).Create(tag).Error; err != nil {
		return exception.InternalError(fmt.Sprintf("Failed to create tag: %v", err))
	}

	return nil
}

// Update modifies an existing tag
func (r *tagRepository) Update(ctx context.Context, tag *entity.Tag) error {
	// Check if tag exists
	var existingTag entity.Tag
	if err := r.db.WithContext(ctx).First(&existingTag, tag.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return exception.NotFound("Tag")
		}
		return exception.InternalError(fmt.Sprintf("Failed to find tag: %v", err))
	}

	// Check slug uniqueness if changed
	if tag.Slug != existingTag.Slug {
		var count int64
		if err := r.db.WithContext(ctx).Model(&entity.Tag{}).
			Where("slug = ? AND id != ?", tag.Slug, tag.ID).
			Count(&count).Error; err != nil {
			return exception.InternalError(fmt.Sprintf("Failed to check slug uniqueness: %v", err))
		}

		if count > 0 {
			return exception.InvalidPayload(map[string]string{
				"slug": "Slug already exists",
			})
		}
	}

	// Check name uniqueness if changed
	if tag.Name != existingTag.Name {
		var count int64
		if err := r.db.WithContext(ctx).Model(&entity.Tag{}).
			Where("name = ? AND id != ?", tag.Name, tag.ID).
			Count(&count).Error; err != nil {
			return exception.InternalError(fmt.Sprintf("Failed to check name uniqueness: %v", err))
		}

		if count > 0 {
			return exception.InvalidPayload(map[string]string{
				"name": "Tag name already exists",
			})
		}
	}

	// Update the tag
	result := r.db.WithContext(ctx).Model(tag).Updates(map[string]interface{}{
		"name":        tag.Name,
		"description": tag.Description,
		"slug":        tag.Slug,
		"updated_at":  time.Now(),
	})

	if result.Error != nil {
		return exception.InternalError(fmt.Sprintf("Failed to update tag: %v", result.Error))
	}

	return nil
}

// Delete removes a tag
func (r *tagRepository) Delete(ctx context.Context, id uint) error {
	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return exception.InternalError(fmt.Sprintf("Failed to begin transaction: %v", tx.Error))
	}

	// Remove tag from articles' associations
	if err := tx.Exec("DELETE FROM article_tags WHERE tag_id = ?", id).Error; err != nil {
		tx.Rollback()
		return exception.InternalError(fmt.Sprintf("Failed to remove tag associations: %v", err))
	}

	// Delete the tag
	result := tx.Delete(&entity.Tag{}, id)
	if result.Error != nil {
		tx.Rollback()
		return exception.InternalError(fmt.Sprintf("Failed to delete tag: %v", result.Error))
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return exception.NotFound("Tag")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return exception.InternalError(fmt.Sprintf("Failed to commit transaction: %v", err))
	}

	return nil
}

// CountBySlug counts tags with the given slug, excluding the specified ID if provided
func (r *tagRepository) CountBySlug(ctx context.Context, slug string, excludeID *uint) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entity.Tag{}).Where("slug = ?", slug)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, exception.InternalError(fmt.Sprintf("Failed to count tags by slug: %v", err))
	}

	return count, nil
}
