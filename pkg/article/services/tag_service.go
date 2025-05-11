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

// tagService implements the TagService interface
type tagService struct {
	tagRepo     repository.TagRepository
	articleRepo repository.ArticleRepository
}

// NewTagService creates a new tag service
func NewTagService(
	tagRepo repository.TagRepository,
	articleRepo repository.ArticleRepository,
) service.TagService {
	return &tagService{
		tagRepo:     tagRepo,
		articleRepo: articleRepo,
	}
}

// Get retrieves a single tag by ID
func (s *tagService) Get(ctx context.Context, id uint) (*entity.Tag, error) {
	return s.tagRepo.FindOne(ctx, id)
}

// GetBySlug retrieves a single tag by its slug
func (s *tagService) GetBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	return s.tagRepo.FindBySlug(ctx, slug)
}

// GetAll retrieves multiple tags with filtering and pagination
func (s *tagService) GetAll(ctx context.Context, query entity.TagQuery) ([]*entity.Tag, int64, error) {
	return s.tagRepo.FindAll(ctx, query)
}

// GetByArticle retrieves all tags associated with an article
func (s *tagService) GetByArticle(ctx context.Context, articleID uint) ([]*entity.Tag, error) {
	// First check if article exists
	if _, err := s.articleRepo.FindOne(ctx, articleID); err != nil {
		return nil, err
	}

	return s.tagRepo.FindByArticleID(ctx, articleID)
}

// generateUniqueSlug creates a unique slug by appending a numeric suffix if needed
func (s *tagService) generateUniqueSlug(ctx context.Context, baseSlug string, excludeID *uint) (string, error) {
	// Try the base slug first
	uniqueSlug := baseSlug
	suffix := 1

	for {
		// Check if slug exists using count
		count, err := s.tagRepo.CountBySlug(ctx, uniqueSlug, excludeID)
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

// Create stores a new tag
func (s *tagService) Create(ctx context.Context, dto entity.TagDTO) (*entity.Tag, error) {
	// Validate name
	if strings.TrimSpace(dto.Name) == "" {
		return nil, exception.InvalidPayload(map[string]string{
			"name": "Tag name is required",
		})
	}

	// Generate slug if not provided
	tagSlug := dto.Slug
	if tagSlug == "" {
		tagSlug = slug.Make(dto.Name)
	}

	// Ensure slug is unique
	uniqueSlug, err := s.generateUniqueSlug(ctx, tagSlug, nil)
	if err != nil {
		return nil, err
	}

	// Create tag entity
	tag := &entity.Tag{
		Name:        dto.Name,
		Description: dto.Description,
		Slug:        uniqueSlug,
	}

	// Store tag
	if err := s.tagRepo.Create(ctx, tag); err != nil {
		return nil, err
	}

	// Return created tag
	return s.tagRepo.FindOne(ctx, tag.ID)
}

// Update modifies an existing tag
func (s *tagService) Update(ctx context.Context, id uint, dto entity.TagDTO) (*entity.Tag, error) {
	// Check tag exists
	tag, err := s.tagRepo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate name
	if strings.TrimSpace(dto.Name) == "" {
		return nil, exception.InvalidPayload(map[string]string{
			"name": "Tag name is required",
		})
	}

	// Update fields
	tag.Name = dto.Name
	tag.Description = dto.Description

	// Generate slug if not provided
	if dto.Slug != "" {
		tag.Slug = dto.Slug
	} else if tag.Name != dto.Name {
		// Generate new slug if name changed and slug not specified
		baseSlug := slug.Make(dto.Name)
		uniqueSlug, err := s.generateUniqueSlug(ctx, baseSlug, &id)
		if err != nil {
			return nil, err
		}
		tag.Slug = uniqueSlug
	}

	// Update tag
	if err := s.tagRepo.Update(ctx, tag); err != nil {
		return nil, err
	}

	// Return updated tag
	return s.tagRepo.FindOne(ctx, tag.ID)
}

// Delete removes a tag
func (s *tagService) Delete(ctx context.Context, id uint) error {
	return s.tagRepo.Delete(ctx, id)
}
