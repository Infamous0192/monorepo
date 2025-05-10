package handlers

import (
	"app/pkg/article/domain/entity"
	"app/pkg/article/domain/service"
	"app/pkg/types/http"
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
func NewCategoryHandler(
	categoryService service.CategoryService,
	validation *validation.Validation,
) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		validation:      validation,
	}
}

// RegisterRoutes registers all routes for category handling
func (h *CategoryHandler) RegisterRoutes(app *fiber.App, authMiddleware fiber.Handler) {
	api := app.Group("/api/categories")

	// Public routes (no API key required)
	api.Get("/", h.GetCategories)
	api.Get("/:id", h.GetCategory)
	api.Get("/slug/:slug", h.GetCategoryBySlug)
	api.Get("/:id/hierarchy", h.GetCategoryHierarchy)
	api.Get("/:id/children", h.GetCategoryChildren)

	// Protected routes (API key required)
	protected := api.Use(authMiddleware)
	protected.Post("/", h.CreateCategory)
	protected.Put("/:id", h.UpdateCategory)
	protected.Delete("/:id", h.DeleteCategory)
}

// GetCategories returns a list of categories with pagination and filtering
// @Summary Get all categories
// @Description Get all categories with pagination and filtering options
// @Tags categories
// @Accept json
// @Produce json
// @Param keyword query string false "Search keyword"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} http.GeneralResponse{data=[]entity.Category}
// @Failure 400 {object} validation.ValidationError
// @Failure 500 {object} error
// @Router /categories [get]
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

	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Categories retrieved successfully",
		Data: fiber.Map{
			"metadata": meta,
			"result":   categories,
		},
	})
}

// GetCategory returns a single category by ID
// @Summary Get category by ID
// @Description Get details of a specific category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} http.GeneralResponse{data=entity.Category}
// @Failure 400 {object} validation.ValidationError
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	category, err := h.categoryService.Get(c.Context(), uint(id))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Category retrieved successfully",
		Data:    category,
	})
}

// GetCategoryBySlug returns a single category by its slug
// @Summary Get category by slug
// @Description Get details of a specific category by its slug
// @Tags categories
// @Accept json
// @Produce json
// @Param slug path string true "Category slug"
// @Success 200 {object} http.GeneralResponse{data=entity.Category}
// @Failure 400 {object} validation.ValidationError
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /categories/slug/{slug} [get]
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

	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Category retrieved successfully",
		Data:    category,
	})
}

// GetCategoryHierarchy returns the complete hierarchy of a category
// @Summary Get category hierarchy
// @Description Get the complete hierarchy (ancestors and descendants) of a category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} http.GeneralResponse{data=entity.CategoryHierarchy}
// @Failure 400 {object} validation.ValidationError
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /categories/{id}/hierarchy [get]
func (h *CategoryHandler) GetCategoryHierarchy(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	hierarchy, err := h.categoryService.GetHierarchy(c.Context(), uint(id))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Category hierarchy retrieved successfully",
		Data:    hierarchy,
	})
}

// GetCategoryChildren returns all direct children of a category
// @Summary Get category children
// @Description Get all direct child categories of a specific category
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} http.GeneralResponse{data=[]entity.Category}
// @Failure 400 {object} validation.ValidationError
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /categories/{id}/children [get]
func (h *CategoryHandler) GetCategoryChildren(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	children, err := h.categoryService.GetChildren(c.Context(), uint(id))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Category children retrieved successfully",
		Data:    children,
	})
}

// CreateCategory creates a new category
// @Summary Create a new category
// @Description Create a new category with the provided information
// @Tags categories
// @Accept json
// @Produce json
// @Param category body entity.CategoryDTO true "Category information"
// @Success 201 {object} http.GeneralResponse{data=entity.Category}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var dto entity.CategoryDTO

	if err := h.validation.Body(&dto, c); err != nil {
		return err
	}

	category, err := h.categoryService.Create(c.Context(), dto)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(http.GeneralResponse{
		Status:  fiber.StatusCreated,
		Message: "Category created successfully",
		Data:    category,
	})
}

// UpdateCategory updates an existing category
// @Summary Update an existing category
// @Description Update a category with the provided information
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body entity.CategoryDTO true "Updated category information"
// @Success 200 {object} http.GeneralResponse{data=entity.Category}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /categories/{id} [put]
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

	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Category updated successfully",
		Data:    category,
	})
}

// DeleteCategory deletes a category
// @Summary Delete a category
// @Description Delete a category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} http.GeneralResponse
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	if err := h.categoryService.Delete(c.Context(), uint(id)); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Category deleted successfully",
	})
}
