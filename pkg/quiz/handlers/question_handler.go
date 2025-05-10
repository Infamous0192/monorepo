package handlers

import (
	"app/pkg/exception"
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/middleware"
	"app/pkg/quiz/services/question"
	"app/pkg/types/http"
	"app/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

// QuestionHandler handles HTTP requests related to questions
type QuestionHandler struct {
	questionService question.QuestionService
	validation      *validation.Validation
}

// NewQuestionHandler creates a new question handler
func NewQuestionHandler(questionService question.QuestionService, validation *validation.Validation) *QuestionHandler {
	return &QuestionHandler{
		questionService: questionService,
		validation:      validation,
	}
}

// RegisterRoutes registers all routes for question handling
func (h *QuestionHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	// Public routes (no authentication required)
	app.Get("/questions", h.GetQuestions)
	app.Get("/questions/:id", h.GetQuestion)

	// Protected routes (admin only)
	protected := app.Group("/questions", authMiddleware.RequireAdmin())
	protected.Post("/", h.CreateQuestion)
	protected.Put("/:id", h.UpdateQuestion)
	protected.Delete("/:id", h.DeleteQuestion)
}

// GetQuestions retrieves a list of questions with optional filtering and pagination
// @Summary Get all questions
// @Description Get all questions with pagination and filtering options
// @Tags questions
// @Accept json
// @Produce json
// @Param quizId query int false "Filter by quiz ID"
// @Param keyword query string false "Search keyword"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} http.GeneralResponse{data=pagination.PaginatedResult[entity.Question]}
// @Failure 400 {object} validation.ValidationError
// @Failure 500 {object} http.GeneralResponse
// @Router /questions [get]
func (h *QuestionHandler) GetQuestions(c *fiber.Ctx) error {
	// Parse query parameters
	query := new(entity.QuestionQuery)
	if err := h.validation.Query(query, c); err != nil {
		return err
	}

	// Set default pagination if not provided
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 10
	}

	// Get questions from service
	result, err := h.questionService.FindAll(c.Context(), *query)
	if err != nil {
		return err
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Questions retrieved successfully",
		Data:    result,
	})
}

// GetQuestion retrieves a single question by ID
// @Summary Get question by ID
// @Description Get details of a specific question by its ID
// @Tags questions
// @Accept json
// @Produce json
// @Param id path int true "Question ID"
// @Success 200 {object} http.GeneralResponse{data=entity.Question}
// @Failure 400 {object} validation.ValidationError
// @Failure 404 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Router /questions/{id} [get]
func (h *QuestionHandler) GetQuestion(c *fiber.Ctx) error {
	// Parse question ID from URL
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	// Get question from service
	question, err := h.questionService.FindOne(c.Context(), uint(id))
	if err != nil {
		return err
	}

	// Check if question exists
	if question == nil {
		return exception.NotFound("question")
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Question retrieved successfully",
		Data:    question,
	})
}

// CreateQuestion creates a new question
// @Summary Create a new question
// @Description Create a new question with the provided information
// @Tags questions
// @Accept json
// @Produce json
// @Param question body entity.QuestionDTO true "Question information"
// @Success 201 {object} http.GeneralResponse{data=entity.Question}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Security BearerAuth
// @Router /questions [post]
func (h *QuestionHandler) CreateQuestion(c *fiber.Ctx) error {
	// Parse request body
	questionDTO := new(entity.QuestionDTO)
	if err := h.validation.Body(questionDTO, c); err != nil {
		return err
	}

	// Create question using service
	question, err := h.questionService.Create(c.Context(), *questionDTO)
	if err != nil {
		return err
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(http.GeneralResponse{
		Status:  fiber.StatusCreated,
		Message: "Question created successfully",
		Data:    question,
	})
}

// UpdateQuestion updates an existing question
// @Summary Update an existing question
// @Description Update a question with the provided information
// @Tags questions
// @Accept json
// @Produce json
// @Param id path int true "Question ID"
// @Param question body entity.QuestionDTO true "Updated question information"
// @Success 200 {object} http.GeneralResponse{data=entity.Question}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} http.GeneralResponse
// @Failure 404 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Security BearerAuth
// @Router /questions/{id} [put]
func (h *QuestionHandler) UpdateQuestion(c *fiber.Ctx) error {
	// Parse question ID from URL
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	// Parse request body
	questionDTO := new(entity.QuestionDTO)
	if err := h.validation.Body(questionDTO, c); err != nil {
		return err
	}

	// Update question using service
	question, err := h.questionService.Update(c.Context(), uint(id), *questionDTO)
	if err != nil {
		return err
	}

	// Check if question exists
	if question == nil {
		return exception.NotFound("question")
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Question updated successfully",
		Data:    question,
	})
}

// DeleteQuestion deletes a question
// @Summary Delete a question
// @Description Delete a question by its ID
// @Tags questions
// @Accept json
// @Produce json
// @Param id path int true "Question ID"
// @Success 200 {object} http.GeneralResponse
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} http.GeneralResponse
// @Failure 404 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Security BearerAuth
// @Router /questions/{id} [delete]
func (h *QuestionHandler) DeleteQuestion(c *fiber.Ctx) error {
	// Parse question ID from URL
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	// Delete question using service
	if err := h.questionService.Delete(c.Context(), uint(id)); err != nil {
		return err
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Question deleted successfully",
	})
}
