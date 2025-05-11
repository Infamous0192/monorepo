package handlers

import (
	"app/pkg/article/domain/entity"
	"app/pkg/article/domain/repository"
	"app/pkg/article/domain/service"
	"app/pkg/exception"
	"app/pkg/types/http"
	"app/pkg/types/pagination"
	"app/pkg/validation"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// ArticleHandler handles HTTP requests related to articles
type ArticleHandler struct {
	articleService service.ArticleService
	fileRepo       repository.FileRepository
	validation     *validation.Validation
}

// NewArticleHandler creates a new article handler
func NewArticleHandler(
	articleService service.ArticleService,
	fileRepo repository.FileRepository,
	validation *validation.Validation,
) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
		fileRepo:       fileRepo,
		validation:     validation,
	}
}

// RegisterRoutes registers all routes for article handling
func (h *ArticleHandler) RegisterRoutes(app fiber.Router, authMiddleware fiber.Handler) {
	api := app.Group("/articles")

	// Public routes (no API key required)
	api.Get("/", h.GetArticles)
	api.Get("/:id", h.GetArticle)
	api.Get("/files/:filename", h.ServeFile)

	// Protected routes (API key required)
	protected := api.Use(authMiddleware)
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
// @Success 200 {object} http.GeneralResponse{data=http.PaginatedResponse{result=[]entity.Article}}
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
		Data: pagination.PaginatedResult[entity.Article]{
			Metadata: meta,
			Result:   articles,
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
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Article title"
// @Param content formData string true "Article content"
// @Param slug formData string false "Article slug (optional)"
// @Param publishedAt formData string false "Published date in RFC3339 format (optional)"
// @Param categoryIds formData []int false "Category IDs (optional)"
// @Param tagIds formData []int false "Tag IDs (optional)"
// @Param thumbnail formData file false "Article thumbnail image (optional)"
// @Success 201 {object} http.GeneralResponse{data=entity.Article}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Security BearerAuth
// @Router /articles [post]
func (h *ArticleHandler) CreateArticle(c *fiber.Ctx) error {
	var dto entity.ArticleDTO

	// Parse article data from form field
	if err := h.validation.Body(&dto, c); err != nil {
		return err
	}

	// Handle thumbnail upload if provided
	if file, err := c.FormFile("thumbnail"); err == nil && file != nil {
		// Validate file type
		if !isValidImageType(file.Header.Get("Content-Type")) {
			return validation.ValidationError{
				Errors: map[string]string{
					"thumbnail": "Invalid file type. Only images are allowed.",
				},
			}
		}

		// Open file
		src, err := file.Open()
		if err != nil {
			return exception.InternalError(fmt.Sprintf("Failed to open uploaded file: %v", err))
		}
		defer src.Close()

		// Store file
		fileEntity, err := h.fileRepo.Store(
			c.Context(),
			file.Filename,
			file.Header.Get("Content-Type"),
			file.Size,
			src,
		)
		if err != nil {
			return err
		}

		// Set thumbnail ID
		dto.ThumbnailID = &fileEntity.ID
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
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Article ID"
// @Param title formData string true "Article title"
// @Param content formData string true "Article content"
// @Param slug formData string false "Article slug (optional)"
// @Param publishedAt formData string false "Published date in RFC3339 format (optional)"
// @Param categoryIds formData []int false "Category IDs (optional)"
// @Param tagIds formData []int false "Tag IDs (optional)"
// @Param thumbnail formData file false "Article thumbnail image (optional)"
// @Success 200 {object} http.GeneralResponse{data=entity.Article}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security BearerAuth
// @Router /articles/{id} [put]
func (h *ArticleHandler) UpdateArticle(c *fiber.Ctx) error {
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	var dto entity.ArticleDTO

	// Parse article data from form field
	if err := h.validation.FormValue(&dto, "article", c); err != nil {
		return err
	}

	// Handle thumbnail upload if provided
	if file, err := c.FormFile("thumbnail"); err == nil && file != nil {
		// Validate file type
		if !isValidImageType(file.Header.Get("Content-Type")) {
			return validation.ValidationError{
				Errors: map[string]string{
					"thumbnail": "Invalid file type. Only images are allowed.",
				},
			}
		}

		// Open file
		src, err := file.Open()
		if err != nil {
			return exception.InternalError(fmt.Sprintf("Failed to open uploaded file: %v", err))
		}
		defer src.Close()

		// Store file
		fileEntity, err := h.fileRepo.Store(
			c.Context(),
			file.Filename,
			file.Header.Get("Content-Type"),
			file.Size,
			src,
		)
		if err != nil {
			return err
		}

		// Set thumbnail ID
		dto.ThumbnailID = &fileEntity.ID
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
// @Security BearerAuth
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

// ServeFile serves uploaded files
func (h *ArticleHandler) ServeFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	if filename == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Filename is required")
	}

	// Get file path from upload directory
	filePath := filepath.Join("uploads", filename) // Make sure this matches your upload directory

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fiber.NewError(fiber.StatusNotFound, "File not found")
	}

	return c.SendFile(filePath)
}

// isValidImageType checks if the content type is a valid image type
func isValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	return validTypes[contentType]
}
