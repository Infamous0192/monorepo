package handlers

import (
	"app/pkg/exception"
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/middleware"
	"app/pkg/quiz/services/quiz"
	"app/pkg/types/http"
	"app/pkg/validation"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// QuizHandler handles HTTP requests related to quizzes
type QuizHandler struct {
	quizService quiz.QuizService
	validation  *validation.Validation
}

// NewQuizHandler creates a new quiz handler
func NewQuizHandler(quizService quiz.QuizService, validation *validation.Validation) *QuizHandler {
	return &QuizHandler{
		quizService: quizService,
		validation:  validation,
	}
}

// RegisterRoutes registers all routes for quiz handling
func (h *QuizHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	// Public routes (no authentication required)
	app.Get("/quizzes", h.GetQuizzes)
	app.Get("/quizzes/:id", h.GetQuiz)

	// Protected routes (admin only)
	protected := app.Group("/quizzes", authMiddleware.RequireAdmin())
	protected.Post("/", h.CreateQuiz)
	protected.Put("/:id", h.UpdateQuiz)
	protected.Delete("/:id", h.DeleteQuiz)
}

// GetQuizzes retrieves a list of quizzes with optional filtering and pagination
// @Summary Get all quizzes
// @Description Get all quizzes with pagination and filtering options
// @Tags quizzes
// @Accept json
// @Produce json
// @Param keyword query string false "Search keyword"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} http.GeneralResponse{data=http.PaginatedResponse{result=[]entity.Quiz}}
// @Failure 400 {object} validation.ValidationError
// @Failure 500 {object} http.GeneralResponse
// @Router /quizzes [get]
func (h *QuizHandler) GetQuizzes(c *fiber.Ctx) error {
	// Parse query parameters
	query := new(entity.QuizQuery)
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

	// Get quizzes from service
	result, err := h.quizService.FindAll(c.Context(), *query)
	if err != nil {
		return err
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Quizzes retrieved successfully",
		Data:    result,
	})
}

// GetQuiz retrieves a single quiz by ID
// @Summary Get quiz by ID
// @Description Get details of a specific quiz by its ID
// @Tags quizzes
// @Accept json
// @Produce json
// @Param id path int true "Quiz ID"
// @Success 200 {object} http.GeneralResponse{data=entity.Quiz}
// @Failure 400 {object} validation.ValidationError
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /quizzes/{id} [get]
func (h *QuizHandler) GetQuiz(c *fiber.Ctx) error {
	// Parse quiz ID from URL
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	// Get quiz from service
	quiz, err := h.quizService.FindOne(c.Context(), uint(id))
	if err != nil {
		return err
	}

	// Check if quiz exists
	if quiz == nil {
		return exception.NotFound("quiz")
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Quiz retrieved successfully",
		Data:    quiz,
	})
}

// CreateQuiz creates a new quiz
// @Summary Create a new quiz
// @Description Create a new quiz with the provided information
// @Tags quizzes
// @Accept json
// @Produce json
// @Param quiz body entity.QuizDTO true "Quiz information"
// @Success 201 {object} http.GeneralResponse{data=entity.Quiz}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Security BearerAuth
// @Router /quizzes [post]
func (h *QuizHandler) CreateQuiz(c *fiber.Ctx) error {
	fmt.Printf("Handler: Received CreateQuiz request with body: %s\n", string(c.Body()))

	// Parse request body
	quizDTO := new(entity.QuizDTO)
	if err := h.validation.Body(quizDTO, c); err != nil {
		fmt.Printf("Handler: Validation error: %v\n", err)
		return err
	}

	fmt.Printf("Handler: Validation passed, creating quiz with DTO: %+v\n", quizDTO)

	// Create quiz using service
	quiz, err := h.quizService.Create(c.Context(), *quizDTO)
	if err != nil {
		fmt.Printf("Handler: Error from service layer: %v\n", err)
		return exception.InternalError(fmt.Sprintf("Failed to create quiz: %v", err))
	}

	fmt.Printf("Handler: Quiz created successfully: %+v\n", quiz)

	// Return response with explicit content type
	c.Set("Content-Type", "application/json")
	response := http.GeneralResponse{
		Status:  fiber.StatusCreated,
		Message: "Quiz created successfully",
		Data:    quiz,
	}
	fmt.Printf("Handler: Sending response: %+v\n", response)

	if err := c.Status(fiber.StatusCreated).JSON(response); err != nil {
		fmt.Printf("Handler: Error sending response: %v\n", err)
		return exception.InternalError(fmt.Sprintf("Failed to send response: %v", err))
	}

	return nil
}

// UpdateQuiz updates an existing quiz
// @Summary Update an existing quiz
// @Description Update a quiz with the provided information
// @Tags quizzes
// @Accept json
// @Produce json
// @Param id path int true "Quiz ID"
// @Param quiz body entity.QuizDTO true "Updated quiz information"
// @Success 200 {object} http.GeneralResponse{data=entity.Quiz}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} http.GeneralResponse
// @Failure 404 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Security BearerAuth
// @Router /quizzes/{id} [put]
func (h *QuizHandler) UpdateQuiz(c *fiber.Ctx) error {
	// Parse quiz ID from URL
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	// Parse request body
	quizDTO := new(entity.QuizDTO)
	if err := h.validation.Body(quizDTO, c); err != nil {
		return err
	}

	// Update quiz using service
	quiz, err := h.quizService.Update(c.Context(), uint(id), *quizDTO)
	if err != nil {
		return err
	}

	// Check if quiz exists
	if quiz == nil {
		return exception.NotFound("quiz")
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Quiz updated successfully",
		Data:    quiz,
	})
}

// DeleteQuiz deletes a quiz
// @Summary Delete a quiz
// @Description Delete a quiz by its ID
// @Tags quizzes
// @Accept json
// @Produce json
// @Param id path int true "Quiz ID"
// @Success 200 {object} http.GeneralResponse
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} http.GeneralResponse
// @Failure 404 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Security BearerAuth
// @Router /quizzes/{id} [delete]
func (h *QuizHandler) DeleteQuiz(c *fiber.Ctx) error {
	// Parse quiz ID from URL
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	// Delete quiz using service
	if err := h.quizService.Delete(c.Context(), uint(id)); err != nil {
		return err
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Quiz deleted successfully",
	})
}
