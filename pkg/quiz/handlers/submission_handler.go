package handlers

import (
	"app/pkg/exception"
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/middleware"
	"app/pkg/quiz/services/submission"
	"app/pkg/types/http"
	"app/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

// SubmissionHandler handles HTTP requests related to submissions
type SubmissionHandler struct {
	submissionService submission.SubmissionService
	validation        *validation.Validation
}

// NewSubmissionHandler creates a new submission handler
func NewSubmissionHandler(submissionService submission.SubmissionService, validation *validation.Validation) *SubmissionHandler {
	return &SubmissionHandler{
		submissionService: submissionService,
		validation:        validation,
	}
}

// RegisterRoutes registers all routes for submission handling
func (h *SubmissionHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	// Authenticated user routes
	authenticated := app.Group("/submissions", authMiddleware.RequireAuth())
	authenticated.Get("/", h.GetSubmissions)
	authenticated.Get("/:id", h.GetSubmission)
	authenticated.Post("/", h.CreateSubmission)
	authenticated.Post("/bulk", h.CreateBulkSubmissions)

	// Admin only routes
	admin := app.Group("/submissions", authMiddleware.RequireAdmin())
	admin.Put("/:id", h.UpdateSubmission)
	admin.Delete("/:id", h.DeleteSubmission)
}

// GetSubmissions retrieves a list of submissions with optional filtering and pagination
// @Summary Get all submissions
// @Description Get all submissions with pagination and filtering options
// @Tags submissions
// @Accept json
// @Produce json
// @Param quizId query int false "Filter by quiz ID"
// @Param userId query int false "Filter by user ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} http.GeneralResponse{data=pagination.PaginatedResult[entity.Submission]}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Security BearerAuth
// @Router /submissions [get]
func (h *SubmissionHandler) GetSubmissions(c *fiber.Ctx) error {
	// Parse query parameters
	query := new(entity.SubmissionQuery)
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

	// Get submissions from service
	result, err := h.submissionService.FindAll(c.Context(), *query)
	if err != nil {
		return err
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Submissions retrieved successfully",
		Data:    result,
	})
}

// GetSubmission retrieves a single submission by ID
// @Summary Get submission by ID
// @Description Get details of a specific submission by its ID
// @Tags submissions
// @Accept json
// @Produce json
// @Param id path int true "Submission ID"
// @Success 200 {object} http.GeneralResponse{data=entity.Submission}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security BearerAuth
// @Router /submissions/{id} [get]
func (h *SubmissionHandler) GetSubmission(c *fiber.Ctx) error {
	// Parse submission ID from URL
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	// Get submission from service
	submission, err := h.submissionService.FindOne(c.Context(), uint(id))
	if err != nil {
		return err
	}

	// Check if submission exists
	if submission == nil {
		return exception.NotFound("submission")
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Submission retrieved successfully",
		Data:    submission,
	})
}

// CreateSubmission creates a new submission
// @Summary Create a new submission
// @Description Create a new submission with the provided information
// @Tags submissions
// @Accept json
// @Produce json
// @Param submission body entity.SubmissionDTO true "Submission information"
// @Success 201 {object} http.GeneralResponse{data=entity.Submission}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Security BearerAuth
// @Router /submissions [post]
func (h *SubmissionHandler) CreateSubmission(c *fiber.Ctx) error {
	// Parse request body
	submissionDTO := new(entity.SubmissionDTO)
	if err := h.validation.Body(submissionDTO, c); err != nil {
		return err
	}

	// Create submission using service
	submission, err := h.submissionService.Create(c.Context(), *submissionDTO)
	if err != nil {
		return err
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(http.GeneralResponse{
		Status:  fiber.StatusCreated,
		Message: "Submission created successfully",
		Data:    submission,
	})
}

// CreateBulkSubmissions creates multiple submissions in a single request
// @Summary Create multiple submissions
// @Description Create multiple submissions in a single request
// @Tags submissions
// @Accept json
// @Produce json
// @Param submissions body entity.SubmissionInsertDTO true "Multiple submission information"
// @Success 201 {object} http.GeneralResponse{data=[]entity.Submission}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Security BearerAuth
// @Router /submissions/bulk [post]
func (h *SubmissionHandler) CreateBulkSubmissions(c *fiber.Ctx) error {
	// Parse request body
	submissionInsertDTO := new(entity.SubmissionInsertDTO)
	if err := h.validation.Body(submissionInsertDTO, c); err != nil {
		return err
	}

	// Check if any submissions were provided
	if len(submissionInsertDTO.Data) == 0 {
		return exception.BadRequest("No submissions provided")
	}

	// Create submissions using service
	submissions, err := h.submissionService.CreateBulk(c.Context(), *submissionInsertDTO)
	if err != nil {
		return err
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(http.GeneralResponse{
		Status:  fiber.StatusCreated,
		Message: "Submissions created successfully",
		Data:    submissions,
	})
}

// UpdateSubmission updates an existing submission
// @Summary Update an existing submission
// @Description Update a submission with the provided information
// @Tags submissions
// @Accept json
// @Produce json
// @Param id path int true "Submission ID"
// @Param submission body entity.SubmissionDTO true "Updated submission information"
// @Success 200 {object} http.GeneralResponse{data=entity.Submission}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security BearerAuth
// @Router /submissions/{id} [put]
func (h *SubmissionHandler) UpdateSubmission(c *fiber.Ctx) error {
	// Parse submission ID from URL
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	// Parse request body
	submissionDTO := new(entity.SubmissionDTO)
	if err := h.validation.Body(submissionDTO, c); err != nil {
		return err
	}

	// Update submission using service
	submission, err := h.submissionService.Update(c.Context(), uint(id), *submissionDTO)
	if err != nil {
		return err
	}

	// Check if submission exists
	if submission == nil {
		return exception.NotFound("submission")
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Submission updated successfully",
		Data:    submission,
	})
}

// DeleteSubmission deletes a submission
// @Summary Delete a submission
// @Description Delete a submission by its ID
// @Tags submissions
// @Accept json
// @Produce json
// @Param id path int true "Submission ID"
// @Success 200 {object} http.GeneralResponse
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security BearerAuth
// @Router /submissions/{id} [delete]
func (h *SubmissionHandler) DeleteSubmission(c *fiber.Ctx) error {
	// Parse submission ID from URL
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	// Delete submission using service
	if err := h.submissionService.Delete(c.Context(), uint(id)); err != nil {
		return err
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Submission deleted successfully",
	})
}
