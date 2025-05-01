package handler

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/chat/service/user"
	"app/pkg/chat/transport/http/dto"
	"app/pkg/chat/transport/http/middleware"
	"app/pkg/exception"
	"app/pkg/types/http"
	"app/pkg/types/pagination"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService      user.UserService
	clientMiddleware *middleware.ClientMiddleware
	authMiddleware   *middleware.AuthMiddleware
}

func NewUserHandler(userService user.UserService, clientMiddleware *middleware.ClientMiddleware, authMiddleware *middleware.AuthMiddleware) *UserHandler {
	return &UserHandler{
		userService:      userService,
		clientMiddleware: clientMiddleware,
		authMiddleware:   authMiddleware,
	}
}

// RegisterRoutes registers all routes for user management
func (h *UserHandler) RegisterRoutes(app fiber.Router) {
	v1 := app.Group("/v1")

	// Protected user routes (requires client key and user authentication)
	users := v1.Group("/users", h.clientMiddleware.ValidateKey(), h.authMiddleware.Authenticate())
	users.Get("/me", h.GetCurrentUser) // Get current user

	// Admin protected routes
	adminUsers := v1.Group("/admin/users", h.clientMiddleware.ValidateKey(), h.authMiddleware.Authenticate())
	adminUsers.Get("/", h.GetUsers)         // List users
	adminUsers.Get("/:id", h.GetUser)       // Get single user
	adminUsers.Post("/", h.CreateUser)      // Create user
	adminUsers.Put("/:id", h.UpdateUser)    // Update any user
	adminUsers.Delete("/:id", h.DeleteUser) // Delete user
}

// GetUsers godoc
// @Summary Get users
// @Description Retrieves users with filtering and pagination
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} http.GeneralResponse{data=http.PaginatedResponse{result=[]entity.User}}
// @Failure 401 {object} http.ErrorResponse
// @Router /v1/users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pag := pagination.Pagination{
		Page:  page,
		Limit: limit,
	}

	users, total, err := h.userService.GetAll(c.Context(), pag)
	if err != nil {
		return err
	}

	metadata := pagination.Metadata{
		Pagination: pag,
		Total:      total,
		Count:      len(users),
		HasPrev:    page > 1,
		HasNext:    len(users) > 0 && int64(page*limit) < total,
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Users fetched successfully",
		Data: map[string]interface{}{
			"metadata": metadata,
			"result":   users,
		},
	})
}

// GetUser godoc
// @Summary Get a user
// @Description Retrieves a single user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} http.GeneralResponse{data=entity.User}
// @Failure 404 {object} http.ErrorResponse
// @Router /v1/users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userService.Get(c.Context(), id)
	if err != nil {
		return err
	}
	if user == nil {
		return exception.NotFound("User")
	}

	return c.JSON(http.GeneralResponse{
		Status: fiber.StatusOK,
		Data:   user,
	})
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Retrieves the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Success 200 {object} http.GeneralResponse{data=entity.User}
// @Failure 401 {object} http.ErrorResponse
// @Router /v1/users/me [get]
func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*entity.User)
	return c.JSON(http.GeneralResponse{
		Status: fiber.StatusOK,
		Data:   user,
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param user body dto.CreateUserRequest true "User details"
// @Success 201 {object} http.GeneralResponse{data=entity.User}
// @Failure 400 {object} http.ErrorResponse
// @Router /v1/admin/users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return exception.BadRequest("Invalid request body")
	}

	user := &entity.User{
		UserID:   req.UserID,
		Name:     req.Name,
		Username: req.Username,
		Picture:  req.Picture,
		Level:    req.Level,
	}

	if err := h.userService.Create(c.Context(), user); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(http.GeneralResponse{
		Status:  fiber.StatusCreated,
		Message: "User created successfully",
		Data:    user,
	})
}

// UpdateUser godoc
// @Summary Update a user
// @Description Updates an existing user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param user body dto.UpdateUserRequest true "User details"
// @Success 200 {object} http.GeneralResponse{data=entity.User}
// @Failure 400,404 {object} http.ErrorResponse
// @Router /v1/users/me [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	var userID string
	currentUser := c.Locals("user").(*entity.User)

	// Check if updating current user or admin updating another user
	if c.Path() == "/v1/users/me" {
		userID = currentUser.ID
	} else {
		userID = c.Params("id")
	}

	// Check if user exists
	user, err := h.userService.Get(c.Context(), userID)
	if err != nil {
		return err
	}
	if user == nil {
		return exception.NotFound("User")
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return exception.BadRequest("Invalid request body")
	}

	// Update fields
	user.Name = req.Name
	user.Username = req.Username
	user.Picture = req.Picture
	user.Level = req.Level

	if err := h.userService.Update(c.Context(), user); err != nil {
		return err
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "User updated successfully",
		Data:    user,
	})
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Deletes an existing user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} http.GeneralResponse{data=entity.User}
// @Failure 404 {object} http.ErrorResponse
// @Router /v1/admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if user exists
	user, err := h.userService.Get(c.Context(), id)
	if err != nil {
		return err
	}
	if user == nil {
		return exception.NotFound("User")
	}

	if err := h.userService.Delete(c.Context(), id); err != nil {
		return err
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "User deleted successfully",
	})
}
