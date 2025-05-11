package handlers

import (
	"app/pkg/exception"
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/middleware"
	"app/pkg/quiz/services/user"
	"app/pkg/types/http"
	"app/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

// UserHandler handles HTTP requests related to user management
type UserHandler struct {
	userService user.UserService
	validation  *validation.Validation
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService user.UserService, validation *validation.Validation) *UserHandler {
	return &UserHandler{
		userService: userService,
		validation:  validation,
	}
}

// RegisterRoutes registers all routes for user handling
func (h *UserHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	// All user routes require admin access
	users := app.Group("/users", authMiddleware.RequireAdmin())
	users.Get("/", h.GetUsers)
	users.Get("/:id", h.GetUser)
	users.Post("/", h.CreateUser)
	users.Put("/:id", h.UpdateUser)
	users.Delete("/:id", h.DeleteUser)
}

// GetUsers retrieves a list of users with optional filtering and pagination
// @Summary Get all users
// @Description Get all users with pagination and filtering options
// @Tags users
// @Accept json
// @Produce json
// @Param email query string false "Filter by email"
// @Param role query string false "Filter by role"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} http.GeneralResponse{data=http.PaginatedResponse{result=[]entity.User}}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 500 {object} error
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	// Parse query parameters
	query := new(entity.UserQuery)
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

	// Get users from service
	result, err := h.userService.FindAll(c.Context(), *query)
	if err != nil {
		return err
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Users retrieved successfully",
		Data:    result,
	})
}

// GetUser retrieves a single user by ID
// @Summary Get user by ID
// @Description Get details of a specific user by its ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} http.GeneralResponse{data=entity.User}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	// Parse user ID from URL
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	// Get user from service
	user, err := h.userService.FindOne(c.Context(), uint(id))
	if err != nil {
		return err
	}

	// Check if user exists
	if user == nil {
		return exception.NotFound("user")
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// CreateUser creates a new user
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body entity.UserDTO true "User information"
// @Success 201 {object} http.GeneralResponse{data=entity.User}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 409 {object} error
// @Failure 500 {object} error
// @Security BearerAuth
// @Router /users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// Parse request body
	userDTO := new(entity.UserDTO)
	if err := h.validation.Body(userDTO, c); err != nil {
		return err
	}

	// Create user using user service
	user, err := h.userService.Create(c.Context(), *userDTO)
	if err != nil {
		return err
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(http.GeneralResponse{
		Status:  fiber.StatusCreated,
		Message: "User created successfully",
		Data:    user,
	})
}

// UpdateUser updates an existing user
// @Summary Update an existing user
// @Description Update a user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body entity.UserDTO true "Updated user information"
// @Success 200 {object} http.GeneralResponse{data=entity.User}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security BearerAuth
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	// Parse user ID from URL
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	// Parse request body
	userDTO := new(entity.UserDTO)
	if err := h.validation.Body(userDTO, c); err != nil {
		return err
	}

	// Update user
	user, err := h.userService.Update(c.Context(), uint(id), *userDTO)
	if err != nil {
		return err
	}

	// Check if user exists
	if user == nil {
		return exception.NotFound("user")
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "User updated successfully",
		Data:    user,
	})
}

// DeleteUser deletes a user
// @Summary Delete a user
// @Description Delete a user by its ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} http.GeneralResponse
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Security BearerAuth
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	// Parse user ID from URL
	id, err := h.validation.ParamsInt(c)
	if err != nil {
		return err
	}

	// Delete user
	if err := h.userService.Delete(c.Context(), uint(id)); err != nil {
		return err
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "User deleted successfully",
	})
}
