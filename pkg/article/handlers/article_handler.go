package handlers

import (
	"app/pkg/article/domain/entity"
	"app/pkg/article/domain/service"
	"app/pkg/types/pagination"
	"app/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

// ArticleHandler handles HTTP requests related to articles
type ArticleHandler struct {
	articleService  service.ArticleService
	categoryService service.CategoryService
	tagService      service.TagService
	validation      *validation.Validation
}

// NewArticleHandler creates a new article handler
func NewArticleHandler(
	articleService service.ArticleService,
	categoryService service.CategoryService,
	tagService service.TagService,
) *ArticleHandler {
	return &ArticleHandler{
		articleService:  articleService,
		categoryService: categoryService,
		tagService:      tagService,
		validation:      validation.NewValidation(),
	}
}

// RegisterRoutes registers all routes for article handling
func (h *ArticleHandler) RegisterRoutes(app *fiber.App, apiKeyMiddleware fiber.Handler) {
	api := app.Group("/api/articles")

	// Public routes (no API key required)
	api.Get("/", h.GetArticles)
	api.Get("/:id", h.GetArticle)
	api.Get("/slug/:slug", h.GetArticleBySlug)

	// Protected routes (API key required)
	protected := api.Use(apiKeyMiddleware)
	protected.Post("/", h.CreateArticle)
	protected.Put("/:id", h.UpdateArticle)
	protected.Delete("/:id", h.DeleteArticle)
	protected.Post("/:id/publish", h.PublishArticle)
	protected.Post("/:id/unpublish", h.UnpublishArticle)
}

// GetArticles returns a list of articles with pagination and filtering
// @Summary Get all articles
// @Description Get all articles with pagination and filtering options
// @Tags articles
// @Accept json
// @Produce json
// @Param keyword query string false "Search keyword"
// @Param categoryId query int false "Filter by category ID"
// @Param published query boolean false "Filter by published status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} fiber.Map
// @Failure 400 {object} validation.ValidationError
// @Failure 500 {object} error
// @Router /articles [get]
func (h *ArticleHandler) GetArticles(c *fiber.Ctx) error {
	query := entity.ArticleQuery{
		Pagination: pagination.Pagination{
			Page:  1,
			Limit: 10,
		},
	}

	// Parse and validate query parameters
	if err := h.validation.Query(&query, c); err != nil {
		return err
	}

	// Handle published filter if provided
	if publishedStr := c.Query("published"); publishedStr != "" {
		published := publishedStr == "true"
		query.IsPublished = &published
	}

	articles, total, err := h.articleService.GetAll(c.Context(), query)
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
		Count:   len(articles),
		HasPrev: query.Page > 1,
		HasNext: int64(query.Page*query.Limit) < total,
	}

	return c.JSON(fiber.Map{
		"status": 200,
		"data": fiber.Map{
			"metadata": meta,
			"result":   articles,
		},
	})
}

// GetArticle returns a single article by ID
// @Summary Get article by ID
// @Description Get details of a specific article by its ID
// @Tags articles
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} validation.ValidationError
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /articles/{id} [get]
func (h *ArticleHandler) GetArticle(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	article, err := h.articleService.Get(c.Context(), uint(id))
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": 200,
		"data":   article,
	})
}

// GetArticleBySlug returns a single article by its slug
// @Summary Get article by slug
// @Description Get details of a specific article by its slug
// @Tags articles
// @Accept json
// @Produce json
// @Param slug path string true "Article slug"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} validation.ValidationError
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /articles/slug/{slug} [get]
func (h *ArticleHandler) GetArticleBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if err := h.validation.Field(slug, "required"); err != nil {
		return validation.ValidationError{
			Errors: map[string]string{
				"slug": err.Error(),
			},
		}
	}

	article, err := h.articleService.GetBySlug(c.Context(), slug)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": 200,
		"data":   article,
	})
}

// CreateArticle creates a new article
// @Summary Create a new article
// @Description Create a new article with the provided information
// @Tags articles
// @Accept json
// @Produce json
// @Param article body entity.ArticleDTO true "Article information"
// @Success 201 {object} fiber.Map
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /articles [post]
func (h *ArticleHandler) CreateArticle(c *fiber.Ctx) error {
	var dto entity.ArticleDTO

	if err := h.validation.Body(&dto, c); err != nil {
		return err
	}

	article, err := h.articleService.Create(c.Context(), dto)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(fiber.Map{
		"status": 201,
		"data":   article,
	})
}

// UpdateArticle updates an existing article
// @Summary Update an existing article
// @Description Update an article with the provided information
// @Tags articles
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Param article body entity.ArticleDTO true "Updated article information"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /articles/{id} [put]
func (h *ArticleHandler) UpdateArticle(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	var dto entity.ArticleDTO
	if err := h.validation.Body(&dto, c); err != nil {
		return err
	}

	article, err := h.articleService.Update(c.Context(), uint(id), dto)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": 200,
		"data":   article,
	})
}

// DeleteArticle deletes an article
// @Summary Delete an article
// @Description Delete an article by its ID
// @Tags articles
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /articles/{id} [delete]
func (h *ArticleHandler) DeleteArticle(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	if err := h.articleService.Delete(c.Context(), uint(id)); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Article deleted successfully",
	})
}

// PublishArticle publishes an article
// @Summary Publish an article
// @Description Change the status of an article to published
// @Tags articles
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /articles/{id}/publish [post]
func (h *ArticleHandler) PublishArticle(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	if err := h.articleService.Publish(c.Context(), uint(id)); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Article published successfully",
	})
}

// UnpublishArticle unpublishes an article
// @Summary Unpublish an article
// @Description Change the status of an article to unpublished
// @Tags articles
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security ApiKeyAuth
// @Router /articles/{id}/unpublish [post]
func (h *ArticleHandler) UnpublishArticle(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	if err := h.articleService.Unpublish(c.Context(), uint(id)); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status":  200,
		"message": "Article unpublished successfully",
	})
}
