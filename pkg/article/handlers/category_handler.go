package handlers

import (
	"app/pkg/article/domain/entity"
	"app/pkg/article/domain/service"
	"app/pkg/types/pagination"
	"app/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

// CategoryHandler handles HTTP requests related to categories
type CategoryHandler struct {
	categoryService service.CategoryService
	validation      *validation.Validation
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		validation:      validation.NewValidation(),
	}
}

// RegisterRoutes registers all routes for category handling
func (h *CategoryHandler) RegisterRoutes(app *fiber.App, apiKeyMiddleware fiber.Handler) {
	api := app.Group("/api/categories")

	// Public routes (no API key required)
	api.Get("/", h.GetCategories)
	api.Get("/:id", h.GetCategory)
	api.Get("/slug/:slug", h.GetCategoryBySlug)
	api.Get("/:id/hierarchy", h.GetCategoryHierarchy)
	api.Get("/:id/children", h.GetCategoryChildren)

	// Protected routes (API key required)
	protected := api.Use(apiKeyMiddleware)
	protected.Post("/", h.CreateCategory)
	protected.Put("/:id", h.UpdateCategory)
	protected.Delete("/:id", h.DeleteCategory)
}

// GetCategories returns a list of categories with pagination and filtering
func (h *CategoryHandler) GetCategories(c *fiber.Ctx) error {
	query := entity.CategoryQuery{
		Pagination: pagination.Pagination{
			Page:  1,
			Limit: 10,
		},
	}

	// Parse and validate query parameters
	if err := h.validation.Query(&query, c); err != nil {
		return err
	}

	categories, total, err := h.categoryService.GetAll(c.Context(), query)
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
		Count:   len(categories),
		HasPrev: query.Page > 1,
		HasNext: int64(query.Page*query.Limit) < total,
	}

	return c.JSON(fiber.Map{
		"status": 200,
		"data": fiber.Map{
			"metadata": meta,
			"result":   categories,
		},
	})
}

// GetCategory returns a single category by ID
func (h *CategoryHandler) GetCategory(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	category, err := h.categoryService.Get(c.Context(), uint(id))
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": 200,
		"data":   category,
	})
}

// GetCategoryBySlug returns a single category by its slug
func (h *CategoryHandler) GetCategoryBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if err := h.validation.Field(slug, "required"); err != nil {
		return validation.ValidationError{
			Errors: map[string]string{
				"slug": err.Error(),
			},
		}
	}

	category, err := h.categoryService.GetBySlug(c.Context(), slug)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": 200,
		"data":   category,
	})
}

// GetCategoryHierarchy returns the complete hierarchy of a category
func (h *CategoryHandler) GetCategoryHierarchy(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	hierarchy, err := h.categoryService.GetHierarchy(c.Context(), uint(id))
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": 200,
		"data":   hierarchy,
	})
}

// GetCategoryChildren returns all direct children of a category
func (h *CategoryHandler) GetCategoryChildren(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	children, err := h.categoryService.GetChildren(c.Context(), uint(id))
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": 200,
		"data":   children,
	})
}

// CreateCategory creates a new category
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var dto entity.CategoryDTO

	if err := h.validation.Body(&dto, c); err != nil {
		return err
	}

	category, err := h.categoryService.Create(c.Context(), dto)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(fiber.Map{
		"status": 201,
		"data":   category,
	})
}

// UpdateCategory updates an existing category
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	var dto entity.CategoryDTO
	if err := h.validation.Body(&dto, c); err != nil {
		return err
	}

	category, err := h.categoryService.Update(c.Context(), uint(id), dto)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": 200,
		"data":   category,
	})
}

// DeleteCategory deletes a category
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	if err := h.categoryService.Delete(c.Context(), uint(id)); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Category deleted successfully",
	})
}
