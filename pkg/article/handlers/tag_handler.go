package handlers

import (
	"app/pkg/article/domain/entity"
	"app/pkg/article/domain/service"
	"app/pkg/types/http"
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
func NewTagHandler(
	tagService service.TagService,
	validation *validation.Validation,
) *TagHandler {
	return &TagHandler{
		tagService: tagService,
		validation: validation,
	}
}

// RegisterRoutes registers all routes for tag handling
func (h *TagHandler) RegisterRoutes(app *fiber.App, authMiddleware fiber.Handler) {
	api := app.Group("/api/tags")

	// Public routes (no API key required)
	api.Get("/", h.GetTags)
	api.Get("/:id", h.GetTag)
	api.Get("/slug/:slug", h.GetTagBySlug)
	api.Get("/article/:articleId", h.GetTagsByArticle)

	// Protected routes (API key required)
	protected := api.Use(authMiddleware)
	protected.Post("/", h.CreateTag)
	protected.Put("/:id", h.UpdateTag)
	protected.Delete("/:id", h.DeleteTag)
}

// GetTags returns a list of tags with pagination and filtering
// @Summary Get all tags
// @Description Get all tags with pagination and filtering options
// @Tags tags
// @Accept json
// @Produce json
// @Param keyword query string false "Search keyword"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} http.GeneralResponse{data=[]entity.Tag}
// @Failure 400 {object} validation.ValidationError
// @Failure 500 {object} error
// @Router /tags [get]
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

	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Tags retrieved successfully",
		Data: fiber.Map{
			"metadata": meta,
			"result":   tags,
		},
	})
}

// GetTag returns a single tag by ID
// @Summary Get tag by ID
// @Description Get details of a specific tag by its ID
// @Tags tags
// @Accept json
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object} http.GeneralResponse{data=entity.Tag}
// @Failure 400 {object} validation.ValidationError
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /tags/{id} [get]
func (h *TagHandler) GetTag(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	tag, err := h.tagService.Get(c.Context(), uint(id))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Tag retrieved successfully",
		Data:    tag,
	})
}

// GetTagBySlug returns a single tag by its slug
// @Summary Get tag by slug
// @Description Get details of a specific tag by its slug
// @Tags tags
// @Accept json
// @Produce json
// @Param slug path string true "Tag slug"
// @Success 200 {object} http.GeneralResponse{data=entity.Tag}
// @Failure 400 {object} validation.ValidationError
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /tags/slug/{slug} [get]
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

	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Tag retrieved successfully",
		Data:    tag,
	})
}

// GetTagsByArticle returns all tags associated with an article
// @Summary Get tags by article
// @Description Get all tags associated with a specific article
// @Tags tags
// @Accept json
// @Produce json
// @Param articleId path int true "Article ID"
// @Success 200 {object} http.GeneralResponse{data=[]entity.Tag}
// @Failure 400 {object} validation.ValidationError
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /tags/article/{articleId} [get]
func (h *TagHandler) GetTagsByArticle(c *fiber.Ctx) error {
	articleId, err := h.validation.ParamsInt(c, "articleId")
	if err != nil {
		return err
	}

	tags, err := h.tagService.GetByArticle(c.Context(), uint(articleId))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Article tags retrieved successfully",
		Data:    tags,
	})
}

// CreateTag creates a new tag
// @Summary Create a new tag
// @Description Create a new tag with the provided information
// @Tags tags
// @Accept json
// @Produce json
// @Param tag body entity.TagDTO true "Tag information"
// @Success 201 {object} http.GeneralResponse{data=entity.Tag}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /tags [post]
func (h *TagHandler) CreateTag(c *fiber.Ctx) error {
	var dto entity.TagDTO

	if err := h.validation.Body(&dto, c); err != nil {
		return err
	}

	tag, err := h.tagService.Create(c.Context(), dto)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(http.GeneralResponse{
		Status:  fiber.StatusCreated,
		Message: "Tag created successfully",
		Data:    tag,
	})
}

// UpdateTag updates an existing tag
// @Summary Update an existing tag
// @Description Update a tag with the provided information
// @Tags tags
// @Accept json
// @Produce json
// @Param id path int true "Tag ID"
// @Param tag body entity.TagDTO true "Updated tag information"
// @Success 200 {object} http.GeneralResponse{data=entity.Tag}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /tags/{id} [put]
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

	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Tag updated successfully",
		Data:    tag,
	})
}

// DeleteTag deletes a tag
// @Summary Delete a tag
// @Description Delete a tag by its ID
// @Tags tags
// @Accept json
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object} http.GeneralResponse
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	if err := h.tagService.Delete(c.Context(), uint(id)); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Tag deleted successfully",
	})
}
