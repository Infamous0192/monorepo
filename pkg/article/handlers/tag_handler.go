package handlers

import (
	"app/pkg/article/domain/entity"
	"app/pkg/article/domain/service"
	"app/pkg/types/pagination"
	"app/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

// TagHandler handles HTTP requests related to tags
type TagHandler struct {
	tagService service.TagService
	validation *validation.Validation
}

// NewTagHandler creates a new tag handler
func NewTagHandler(tagService service.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
		validation: validation.NewValidation(),
	}
}

// RegisterRoutes registers all routes for tag handling
func (h *TagHandler) RegisterRoutes(app *fiber.App, apiKeyMiddleware fiber.Handler) {
	api := app.Group("/api/tags")

	// Public routes (no API key required)
	api.Get("/", h.GetTags)
	api.Get("/:id", h.GetTag)
	api.Get("/slug/:slug", h.GetTagBySlug)
	api.Get("/article/:articleId", h.GetTagsByArticle)

	// Protected routes (API key required)
	protected := api.Use(apiKeyMiddleware)
	protected.Post("/", h.CreateTag)
	protected.Put("/:id", h.UpdateTag)
	protected.Delete("/:id", h.DeleteTag)
}

// GetTags returns a list of tags with pagination and filtering
func (h *TagHandler) GetTags(c *fiber.Ctx) error {
	query := entity.TagQuery{
		Pagination: pagination.Pagination{
			Page:  1,
			Limit: 10,
		},
	}

	// Parse and validate query parameters
	if err := h.validation.Query(&query, c); err != nil {
		return err
	}

	tags, total, err := h.tagService.GetAll(c.Context(), query)
	if err != nil {
		return err
	}

	// Build response with pagination metadata
	meta := pagination.Metadata{
		Pagination: pagination.Pagination{
			Page:  query.Page,
			Limit: query.Limit,
		},
		Total:   total,
		Count:   len(tags),
		HasPrev: query.Page > 1,
		HasNext: int64(query.Page*query.Limit) < total,
	}

	return c.JSON(fiber.Map{
		"status": 200,
		"data": fiber.Map{
			"metadata": meta,
			"result":   tags,
		},
	})
}

// GetTag returns a single tag by ID
func (h *TagHandler) GetTag(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	tag, err := h.tagService.Get(c.Context(), uint(id))
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": 200,
		"data":   tag,
	})
}

// GetTagBySlug returns a single tag by its slug
func (h *TagHandler) GetTagBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if err := h.validation.Field(slug, "required"); err != nil {
		return validation.ValidationError{
			Errors: map[string]string{
				"slug": err.Error(),
			},
		}
	}

	tag, err := h.tagService.GetBySlug(c.Context(), slug)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": 200,
		"data":   tag,
	})
}

// GetTagsByArticle returns all tags associated with an article
func (h *TagHandler) GetTagsByArticle(c *fiber.Ctx) error {
	articleId, err := h.validation.ParamsInt(c, "articleId")
	if err != nil {
		return err
	}

	tags, err := h.tagService.GetByArticle(c.Context(), uint(articleId))
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": 200,
		"data":   tags,
	})
}

// CreateTag creates a new tag
func (h *TagHandler) CreateTag(c *fiber.Ctx) error {
	var dto entity.TagDTO

	if err := h.validation.Body(&dto, c); err != nil {
		return err
	}

	tag, err := h.tagService.Create(c.Context(), dto)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(fiber.Map{
		"status": 201,
		"data":   tag,
	})
}

// UpdateTag updates an existing tag
func (h *TagHandler) UpdateTag(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	var dto entity.TagDTO
	if err := h.validation.Body(&dto, c); err != nil {
		return err
	}

	tag, err := h.tagService.Update(c.Context(), uint(id), dto)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": 200,
		"data":   tag,
	})
}

// DeleteTag deletes a tag
func (h *TagHandler) DeleteTag(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	if err := h.tagService.Delete(c.Context(), uint(id)); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Tag deleted successfully",
	})
}
