package services

import (
	"app/pkg/article/domain/entity"
	"app/pkg/article/domain/repository"
	"app/pkg/article/domain/service"
	"app/pkg/exception"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gosimple/slug"
)

// articleService implements the ArticleService interface
type articleService struct {
	articleRepo  repository.ArticleRepository
	categoryRepo repository.CategoryRepository
	tagRepo      repository.TagRepository
	fileRepo     repository.FileRepository
}

// NewArticleService creates a new article service
func NewArticleService(
	articleRepo repository.ArticleRepository,
	categoryRepo repository.CategoryRepository,
	tagRepo repository.TagRepository,
	fileRepo repository.FileRepository,
) service.ArticleService {
	return &articleService{
		articleRepo:  articleRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
		fileRepo:     fileRepo,
	}
}

// Get retrieves a single article by ID
func (s *articleService) Get(ctx context.Context, id uint) (*entity.Article, error) {
	return s.articleRepo.FindOne(ctx, id)
}

// GetBySlug retrieves a single article by its slug
func (s *articleService) GetBySlug(ctx context.Context, slug string) (*entity.Article, error) {
	return s.articleRepo.FindBySlug(ctx, slug)
}

// GetAll retrieves multiple articles with filtering and pagination
func (s *articleService) GetAll(ctx context.Context, query entity.ArticleQuery) ([]*entity.Article, int64, error) {
	return s.articleRepo.FindAll(ctx, query)
}

// generateUniqueSlug creates a unique slug by appending a numeric suffix if needed
func (s *articleService) generateUniqueSlug(ctx context.Context, baseSlug string, excludeID *uint) (string, error) {
	// Try the base slug first
	uniqueSlug := baseSlug
	suffix := 1

	for {
		// Check if slug exists using count
		count, err := s.articleRepo.CountBySlug(ctx, uniqueSlug, excludeID)
		if err != nil {
			return "", err
		}

		// If count is 0 or it's the same article we're updating, slug is unique
		if count == 0 {
			return uniqueSlug, nil
		}

		// Try next suffix
		uniqueSlug = fmt.Sprintf("%s-%d", baseSlug, suffix)
		suffix++
	}
}

// Create stores a new article
func (s *articleService) Create(ctx context.Context, dto entity.ArticleDTO) (*entity.Article, error) {
	// Generate slug if not provided
	articleSlug := dto.Slug
	if articleSlug == "" {
		articleSlug = slug.Make(dto.Title)
	}

	// Ensure slug is unique
	uniqueSlug, err := s.generateUniqueSlug(ctx, articleSlug, nil)
	if err != nil {
		return nil, err
	}

	// Create article entity
	article := &entity.Article{
		Title:       dto.Title,
		Content:     dto.Content,
		Slug:        uniqueSlug,
		ThumbnailID: dto.ThumbnailID,
		PublishedAt: dto.PublishedAt,
	}

	// Process categories if provided
	if len(dto.CategoryIDs) > 0 {
		categories := make([]*entity.Category, 0, len(dto.CategoryIDs))
		for _, categoryID := range dto.CategoryIDs {
			category, err := s.categoryRepo.FindOne(ctx, categoryID)
			if err != nil {
				return nil, exception.InvalidPayload(map[string]string{
					"categoryIds": "Invalid category ID",
				})
			}
			categories = append(categories, category)
		}
		article.Categories = categories
	}

	// Process tags if provided
	if len(dto.TagIDs) > 0 {
		tags := make([]*entity.Tag, 0, len(dto.TagIDs))
		for _, tagID := range dto.TagIDs {
			tag, err := s.tagRepo.FindOne(ctx, tagID)
			if err != nil {
				return nil, exception.InvalidPayload(map[string]string{
					"tagIds": "Invalid tag ID",
				})
			}
			tags = append(tags, tag)
		}
		article.Tags = tags
	}

	// Store article
	if err := s.articleRepo.Create(ctx, article); err != nil {
		return nil, err
	}

	// Return created article with all relationships loaded
	return s.articleRepo.FindOne(ctx, article.ID)
}

// Update modifies an existing article
func (s *articleService) Update(ctx context.Context, id uint, dto entity.ArticleDTO) (*entity.Article, error) {
	// Check article exists
	article, err := s.articleRepo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	article.Title = dto.Title
	article.Content = dto.Content

	// Generate slug if not provided
	if dto.Slug != "" {
		article.Slug = dto.Slug
	} else if article.Title != dto.Title {
		// Generate new slug if title changed and slug not specified
		baseSlug := slug.Make(dto.Title)
		uniqueSlug, err := s.generateUniqueSlug(ctx, baseSlug, &id)
		if err != nil {
			return nil, err
		}
		article.Slug = uniqueSlug
	}

	// Update published status if provided
	if dto.PublishedAt != nil {
		article.PublishedAt = dto.PublishedAt
	}

	// Process categories if provided
	if len(dto.CategoryIDs) > 0 {
		categories := make([]*entity.Category, 0, len(dto.CategoryIDs))
		for _, categoryID := range dto.CategoryIDs {
			category, err := s.categoryRepo.FindOne(ctx, categoryID)
			if err != nil {
				return nil, exception.InvalidPayload(map[string]string{
					"categoryIds": "Invalid category ID",
				})
			}
			categories = append(categories, category)
		}
		article.Categories = categories
	}

	// Process tags if provided
	if len(dto.TagIDs) > 0 {
		tags := make([]*entity.Tag, 0, len(dto.TagIDs))
		for _, tagID := range dto.TagIDs {
			tag, err := s.tagRepo.FindOne(ctx, tagID)
			if err != nil {
				return nil, exception.InvalidPayload(map[string]string{
					"tagIds": "Invalid tag ID",
				})
			}
			tags = append(tags, tag)
		}
		article.Tags = tags
	}

	// Update article
	if err := s.articleRepo.Update(ctx, article); err != nil {
		return nil, err
	}

	// Return updated article with all relationships loaded
	return s.articleRepo.FindOne(ctx, article.ID)
}

// Delete removes an article
func (s *articleService) Delete(ctx context.Context, id uint) error {
	// Get article first to check if it has a thumbnail
	article, err := s.articleRepo.FindOne(ctx, id)
	if err != nil {
		return err
	}

	// Delete thumbnail if exists
	if article.ThumbnailID != nil {
		if err := s.fileRepo.Delete(ctx, *article.ThumbnailID); err != nil {
			return err
		}
	}

	// Delete article
	return s.articleRepo.Delete(ctx, id)
}

// Publish sets an article as published
func (s *articleService) Publish(ctx context.Context, id uint) error {
	// Check article exists
	article, err := s.articleRepo.FindOne(ctx, id)
	if err != nil {
		return err
	}

	// Validate article has content
	if strings.TrimSpace(article.Content) == "" {
		return exception.BadRequest("Cannot publish an article without content")
	}

	// Set published date to now
	now := time.Now().Format(time.RFC3339)
	return s.articleRepo.Publish(ctx, id, &now)
}

// Unpublish removes the published status of an article
func (s *articleService) Unpublish(ctx context.Context, id uint) error {
	// Check article exists
	if _, err := s.articleRepo.FindOne(ctx, id); err != nil {
		return err
	}

	// Set published date to nil
	return s.articleRepo.Publish(ctx, id, nil)
}
