package handlers

import (
	"app/pkg/article/domain/entity"
	"app/pkg/article/domain/service"
	"app/pkg/types/http"
	"app/pkg/types/pagination"
	"app/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

// ArticleHandler handles HTTP requests related to articles
type ArticleHandler struct {
	articleService service.ArticleService
	validation     *validation.Validation
}

// NewArticleHandler creates a new article handler
func NewArticleHandler(
	articleService service.ArticleService,
	validation *validation.Validation,
) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
		validation:     validation,
	}
}

// RegisterRoutes registers all routes for article handling
func (h *ArticleHandler) RegisterRoutes(app *fiber.App, apiKeyMiddleware fiber.Handler) {
	api := app.Group("/api/articles")

	// Public routes (no API key required)
	api.Get("/", h.GetArticles)
	api.Get("/:id", h.GetArticle)

	// Protected routes (API key required)
	protected := api.Use(apiKeyMiddleware)
	protected.Post("/", h.CreateArticle)
	protected.Put("/:id", h.UpdateArticle)
	protected.Delete("/:id", h.DeleteArticle)
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
// @Success 200 {object} http.GeneralResponse{data=[]entity.Article}
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

	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Articles retrieved successfully",
		Data: fiber.Map{
			"metadata": meta,
			"result":   articles,
		},
	})
}

// GetArticle returns a single article by its slug
// @Summary Get article by slug
// @Description Get details of a specific article by its slug
// @Tags articles
// @Accept json
// @Produce json
// @Param slug path string true "Article slug"
// @Success 200 {object} http.GeneralResponse{data=entity.Article}
// @Failure 400 {object} validation.ValidationError
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /articles/{slug} [get]
func (h *ArticleHandler) GetArticle(c *fiber.Ctx) error {
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

	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Article retrieved successfully",
		Data:    article,
	})
}

// CreateArticle creates a new article
// @Summary Create a new article
// @Description Create a new article with the provided information
// @Tags articles
// @Accept json
// @Produce json
// @Param article body entity.ArticleDTO true "Article information"
// @Success 201 {object} http.GeneralResponse{data=entity.Article}
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

	return c.Status(fiber.StatusCreated).JSON(http.GeneralResponse{
		Status:  fiber.StatusCreated,
		Message: "Article created successfully",
		Data:    article,
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
// @Success 200 {object} http.GeneralResponse{data=entity.Article}
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

	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Article updated successfully",
		Data:    article,
	})
}

// DeleteArticle deletes an article
// @Summary Delete an article
// @Description Delete an article by its ID
// @Tags articles
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Success 200 {object} http.GeneralResponse
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

	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Article deleted successfully",
	})
}
