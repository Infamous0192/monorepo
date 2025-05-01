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

// articleRepository implements the ArticleRepository interface using GORM
type articleRepository struct {
	db *gorm.DB
}

// NewArticleRepository creates a new instance of the article repository
func NewArticleRepository(db *gorm.DB) repository.ArticleRepository {
	return &articleRepository{
		db: db,
	}
}

// FindOne retrieves a single article by ID
func (r *articleRepository) FindOne(ctx context.Context, id uint) (*entity.Article, error) {
	var article entity.Article

	result := r.db.WithContext(ctx).
		Preload("Categories").
		Preload("Tags").
		First(&article, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, exception.NotFound("Article")
		}
		return nil, exception.InternalError(fmt.Sprintf("Failed to find article: %v", result.Error))
	}

	return &article, nil
}

// FindBySlug retrieves a single article by its slug
func (r *articleRepository) FindBySlug(ctx context.Context, slug string) (*entity.Article, error) {
	var article entity.Article

	result := r.db.WithContext(ctx).
		Preload("Categories").
		Preload("Tags").
		Where("slug = ?", slug).
		First(&article)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, exception.NotFound("Article")
		}
		return nil, exception.InternalError(fmt.Sprintf("Failed to find article by slug: %v", result.Error))
	}

	return &article, nil
}

// FindAll retrieves multiple articles with filtering and pagination
func (r *articleRepository) FindAll(ctx context.Context, query entity.ArticleQuery) ([]*entity.Article, int64, error) {
	var articles []*entity.Article
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.Article{})

	// Apply filters
	if query.Keyword != "" {
		db = db.Where("title LIKE ? OR content LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	if query.CategoryID != 0 {
		db = db.Joins("JOIN article_categories ON article_categories.article_id = articles.id").
			Where("article_categories.category_id = ?", query.CategoryID)
	}

	if query.TagID != 0 {
		db = db.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id = ?", query.TagID)
	}

	if query.IsPublished != nil {
		if *query.IsPublished {
			db = db.Where("published_at IS NOT NULL")
		} else {
			db = db.Where("published_at IS NULL")
		}
	}

	// Count total records first
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, exception.InternalError(fmt.Sprintf("Failed to count articles: %v", err))
	}

	// Apply pagination and retrieve data
	result := db.
		Preload("Categories").
		Preload("Tags").
		Limit(query.GetLimit()).
		Offset(query.GetOffset()).
		Order("created_at DESC").
		Find(&articles)

	if result.Error != nil {
		return nil, 0, exception.InternalError(fmt.Sprintf("Failed to find articles: %v", result.Error))
	}

	return articles, total, nil
}

// Create stores a new article
func (r *articleRepository) Create(ctx context.Context, article *entity.Article) error {
	// Check if slug already exists
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Article{}).
		Where("slug = ?", article.Slug).
		Count(&count).Error; err != nil {
		return exception.InternalError(fmt.Sprintf("Failed to check slug uniqueness: %v", err))
	}

	if count > 0 {
		return exception.InvalidPayload(map[string]string{
			"slug": "Slug already exists",
		})
	}

	if err := r.db.WithContext(ctx).Create(article).Error; err != nil {
		return exception.InternalError(fmt.Sprintf("Failed to create article: %v", err))
	}

	return nil
}

// Update modifies an existing article
func (r *articleRepository) Update(ctx context.Context, article *entity.Article) error {
	// Check if article exists
	var existingArticle entity.Article
	if err := r.db.WithContext(ctx).First(&existingArticle, article.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return exception.NotFound("Article")
		}
		return exception.InternalError(fmt.Sprintf("Failed to find article: %v", err))
	}

	// Check slug uniqueness if changed
	if article.Slug != existingArticle.Slug {
		var count int64
		if err := r.db.WithContext(ctx).Model(&entity.Article{}).
			Where("slug = ? AND id != ?", article.Slug, article.ID).
			Count(&count).Error; err != nil {
			return exception.InternalError(fmt.Sprintf("Failed to check slug uniqueness: %v", err))
		}

		if count > 0 {
			return exception.InvalidPayload(map[string]string{
				"slug": "Slug already exists",
			})
		}
	}

	// Update the article
	result := r.db.WithContext(ctx).Model(article).Updates(map[string]interface{}{
		"title":        article.Title,
		"content":      article.Content,
		"slug":         article.Slug,
		"published_at": article.PublishedAt,
		"updated_at":   time.Now(),
	})

	if result.Error != nil {
		return exception.InternalError(fmt.Sprintf("Failed to update article: %v", result.Error))
	}

	// Update associations if needed
	if len(article.Categories) > 0 || len(article.Tags) > 0 {
		// Start a transaction
		tx := r.db.WithContext(ctx).Begin()
		if tx.Error != nil {
			return exception.InternalError(fmt.Sprintf("Failed to begin transaction: %v", tx.Error))
		}

		// Update category associations
		if len(article.Categories) > 0 {
			if err := tx.Model(article).Association("Categories").Replace(article.Categories); err != nil {
				tx.Rollback()
				return exception.InternalError(fmt.Sprintf("Failed to update article categories: %v", err))
			}
		}

		// Update tag associations
		if len(article.Tags) > 0 {
			if err := tx.Model(article).Association("Tags").Replace(article.Tags); err != nil {
				tx.Rollback()
				return exception.InternalError(fmt.Sprintf("Failed to update article tags: %v", err))
			}
		}

		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			return exception.InternalError(fmt.Sprintf("Failed to commit transaction: %v", err))
		}
	}

	return nil
}

// Delete removes an article
func (r *articleRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&entity.Article{}, id)

	if result.Error != nil {
		return exception.InternalError(fmt.Sprintf("Failed to delete article: %v", result.Error))
	}

	if result.RowsAffected == 0 {
		return exception.NotFound("Article")
	}

	return nil
}

// Publish sets the published status of an article
func (r *articleRepository) Publish(ctx context.Context, id uint, publishedAt *string) error {
	var timestamp *time.Time
	if publishedAt != nil {
		parsedTime, err := time.Parse(time.RFC3339, *publishedAt)
		if err != nil {
			return exception.InvalidPayload(map[string]string{
				"publishedAt": "Invalid date format. Use RFC3339 format.",
			})
		}
		timestamp = &parsedTime
	} else {
		now := time.Now()
		timestamp = &now
	}

	result := r.db.WithContext(ctx).Model(&entity.Article{}).
		Where("id = ?", id).
		Update("published_at", timestamp)

	if result.Error != nil {
		return exception.InternalError(fmt.Sprintf("Failed to publish article: %v", result.Error))
	}

	if result.RowsAffected == 0 {
		return exception.NotFound("Article")
	}

	return nil
}
