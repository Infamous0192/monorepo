package handlers

import (
	"app/pkg/exception"
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/services/auth"
	"app/pkg/types/http"
	"app/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService auth.AuthService
	validation  *validation.Validation
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validation:  validation.NewValidation(),
	}
}

// Register handles user registration requests
// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags auth
// @Accept json
// @Produce json
// @Param user body entity.UserDTO true "User registration information"
// @Success 201 {object} http.GeneralResponse
// @Failure 400 {object} http.GeneralResponse
// @Failure 409 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// Parse request body
	userDTO := new(entity.UserDTO)
	if err := h.validation.Body(userDTO, c); err != nil {
		return err
	}

	// Register the user
	user, err := h.authService.Register(userDTO)
	if err != nil {
		return err
	}

	// Return the created user
	return c.Status(fiber.StatusCreated).JSON(http.GeneralResponse{
		Status:  fiber.StatusCreated,
		Message: "User registered successfully",
		Data:    user,
	})
}

// Login handles user login requests
// @Summary Login a user
// @Description Login with username and password to get an authentication token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body object true "Login credentials"
// @Success 200 {object} http.GeneralResponse
// @Failure 400 {object} http.GeneralResponse
// @Failure 401 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Parse request body
	loginRequest := new(struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	})

	if err := h.validation.Body(loginRequest, c); err != nil {
		return err
	}

	// Login the user
	token, err := h.authService.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		return err
	}

	// Return the JWT token
	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Login successful",
		Data: fiber.Map{
			"token": token,
		},
	})
}

// GetProfile returns the current user's profile
// @Summary Get user profile
// @Description Get the profile of the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} http.GeneralResponse
// @Failure 401 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Security BearerAuth
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	// Get user from context (set by middleware)
	claims, ok := c.Locals("user").(*auth.Claims)
	if !ok {
		return exception.Http(401, "User not authenticated")
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "User profile retrieved",
		Data:    claims,
	})
}
