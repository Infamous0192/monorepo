package handlers

import (
	"app/pkg/exception"
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/services"
	"app/pkg/types/http"
	"app/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

// AnswerHandler handles HTTP requests related to answers
type AnswerHandler struct {
	services   *services.Services
	validation *validation.Validation
}

// NewAnswerHandler creates a new answer handler
func NewAnswerHandler(services *services.Services) *AnswerHandler {
	return &AnswerHandler{
		services:   services,
		validation: validation.NewValidation(),
	}
}

// GetAnswers retrieves a list of answers with optional filtering and pagination
// @Summary Get all answers
// @Description Get all answers with pagination and filtering options
// @Tags answers
// @Accept json
// @Produce json
// @Param questionId query int false "Filter by question ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} http.GeneralResponse
// @Failure 400 {object} http.GeneralResponse
// @Failure 401 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Security BearerAuth
// @Router /answers [get]
func (h *AnswerHandler) GetAnswers(c *fiber.Ctx) error {
	// Parse query parameters
	query := new(entity.AnswerQuery)
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

	// Get answers from service
	result, err := h.services.Answer.FindAll(c.Context(), *query)
	if err != nil {
		return err
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Answers retrieved successfully",
		Data:    result,
	})
}

// GetAnswer retrieves a single answer by ID
// @Summary Get answer by ID
// @Description Get details of a specific answer by its ID
// @Tags answers
// @Accept json
// @Produce json
// @Param id path int true "Answer ID"
// @Success 200 {object} http.GeneralResponse
// @Failure 400 {object} http.GeneralResponse
// @Failure 401 {object} http.GeneralResponse
// @Failure 404 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Security BearerAuth
// @Router /answers/{id} [get]
func (h *AnswerHandler) GetAnswer(c *fiber.Ctx) error {
	// Parse answer ID from URL
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	// Get answer from service
	answer, err := h.services.Answer.FindOne(c.Context(), uint(id))
	if err != nil {
		return err
	}

	// Check if answer exists
	if answer == nil {
		return exception.NotFound("answer")
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Answer retrieved successfully",
		Data:    answer,
	})
}

// CreateAnswer creates a new answer
// @Summary Create a new answer
// @Description Create a new answer with the provided information
// @Tags answers
// @Accept json
// @Produce json
// @Param answer body entity.AnswerDTO true "Answer information"
// @Success 201 {object} http.GeneralResponse
// @Failure 400 {object} http.GeneralResponse
// @Failure 401 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Security BearerAuth
// @Router /answers [post]
func (h *AnswerHandler) CreateAnswer(c *fiber.Ctx) error {
	// Parse request body
	answerDTO := new(entity.AnswerDTO)
	if err := h.validation.Body(answerDTO, c); err != nil {
		return err
	}

	// Create answer using service
	answer, err := h.services.Answer.Create(c.Context(), *answerDTO)
	if err != nil {
		return err
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(http.GeneralResponse{
		Status:  fiber.StatusCreated,
		Message: "Answer created successfully",
		Data:    answer,
	})
}

// UpdateAnswer updates an existing answer
// @Summary Update an existing answer
// @Description Update an answer with the provided information
// @Tags answers
// @Accept json
// @Produce json
// @Param id path int true "Answer ID"
// @Param answer body entity.AnswerDTO true "Updated answer information"
// @Success 200 {object} http.GeneralResponse
// @Failure 400 {object} http.GeneralResponse
// @Failure 401 {object} http.GeneralResponse
// @Failure 404 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Security BearerAuth
// @Router /answers/{id} [put]
func (h *AnswerHandler) UpdateAnswer(c *fiber.Ctx) error {
	// Parse answer ID from URL
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	// Parse request body
	answerDTO := new(entity.AnswerDTO)
	if err := h.validation.Body(answerDTO, c); err != nil {
		return err
	}

	// Update answer using service
	answer, err := h.services.Answer.Update(c.Context(), uint(id), *answerDTO)
	if err != nil {
		return err
	}

	// Check if answer exists
	if answer == nil {
		return exception.NotFound("answer")
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Answer updated successfully",
		Data:    answer,
	})
}

// DeleteAnswer deletes an answer
// @Summary Delete an answer
// @Description Delete an answer by its ID
// @Tags answers
// @Accept json
// @Produce json
// @Param id path int true "Answer ID"
// @Success 200 {object} http.GeneralResponse
// @Failure 400 {object} http.GeneralResponse
// @Failure 401 {object} http.GeneralResponse
// @Failure 404 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Security BearerAuth
// @Router /answers/{id} [delete]
func (h *AnswerHandler) DeleteAnswer(c *fiber.Ctx) error {
	// Parse answer ID from URL
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	// Delete answer using service
	if err := h.services.Answer.Delete(c.Context(), uint(id)); err != nil {
		return err
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Answer deleted successfully",
	})
}
