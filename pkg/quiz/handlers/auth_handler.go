package handlers

import (
	"app/pkg/exception"
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/middleware"
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
func NewAuthHandler(authService auth.AuthService, validation *validation.Validation) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validation:  validation,
	}
}

// RegisterRoutes registers all routes for authentication handling
func (h *AuthHandler) RegisterRoutes(app fiber.Router) {
	// Public routes (no authentication required)
	auth := app.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)

	// Protected routes (authentication required)
	authMiddleware := middleware.NewAuthMiddleware(h.authService)
	auth.Get("/verify", authMiddleware.RequireAuth(), h.Verify)
}

// Register handles user registration requests
// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags auth
// @Accept json
// @Produce json
// @Param user body entity.RegisterDTO true "User registration information"
// @Success 201 {object} http.GeneralResponse{data=entity.User}
// @Failure 400 {object} validation.ValidationError
// @Failure 409 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	// Parse request body
	registerDTO := new(entity.RegisterDTO)
	if err := h.validation.Body(registerDTO, c); err != nil {
		return err
	}

	// Register the user
	user, err := h.authService.Register(registerDTO)
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
// @Param credentials body entity.LoginDTO true "Login credentials"
// @Success 200 {object} http.GeneralResponse{data=map[string]string}
// @Failure 400 {object} validation.ValidationError
// @Failure 401 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Parse request body
	loginDTO := new(entity.LoginDTO)
	if err := h.validation.Body(loginDTO, c); err != nil {
		return err
	}

	// Login the user
	token, err := h.authService.Login(loginDTO.Username, loginDTO.Password)
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

// Verify verifies the current user's authentication token
// @Summary Verify authentication token
// @Description Verify the current user's authentication token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} http.GeneralResponse{data=auth.Claims}
// @Failure 401 {object} http.GeneralResponse
// @Failure 500 {object} http.GeneralResponse
// @Security BearerAuth
// @Router /auth/verify [get]
func (h *AuthHandler) Verify(c *fiber.Ctx) error {
	// Get user from context (set by middleware)
	claims, ok := c.Locals("user").(*auth.Claims)
	if !ok {
		return exception.Http(401, "User not authenticated")
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Token verified successfully",
		Data:    claims,
	})
}
