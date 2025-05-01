package telegram

import (
	"app/pkg/exception"
	"app/pkg/types/http"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles HTTP requests for Telegram authentication
type AuthHandler struct {
	authService *AuthService
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterRoutes registers all routes for authentication
func (h *AuthHandler) RegisterRoutes(app fiber.Router) {
	auth := app.Group("/auth/telegram")

	auth.Post("/login", h.Login)
}

// Login godoc
// @Summary Login with Telegram
// @Description Validates Telegram Web App initData and returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request with initData"
// @Success 200 {object} http.GeneralResponse{data=LoginResponse}
// @Failure 400,401 {object} http.ErrorResponse
// @Router /auth/telegram/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return exception.BadRequest("Invalid request format")
	}

	if req.InitData == "" {
		return exception.BadRequest("initData is required")
	}

	user, err := h.authService.ValidateInitData(req.InitData)
	if err != nil {
		return err
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(user)
	if err != nil {
		return exception.InternalError("Failed to generate token")
	}

	response := LoginResponse{
		Token: token,
		User:  *user,
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Login successful",
		Data:    response,
	})
}
