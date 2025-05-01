package twitter

import (
	"app/pkg/exception"
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// LoginResponse represents the response for successful login
type LoginResponse struct {
	Token string      `json:"token"`
	User  TwitterUser `json:"user"`
}

type AuthHandler struct {
	service *AuthService
}

// NewAuthHandler creates a new Twitter authentication handler
func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// RegisterRoutes registers the auth routes
func (h *AuthHandler) RegisterRoutes(app *fiber.App) {
	auth := app.Group("/auth/twitter")
	auth.Get("/login", h.Login)
	auth.Get("/callback", h.Callback)
}

// Login initiates Twitter OAuth2 flow
// @Summary Initiate Twitter OAuth2 login
// @Description Redirects to Twitter OAuth2 authorization page
// @Tags auth
// @Produce json
// @Success 302 {string} string "Redirect to Twitter"
// @Router /auth/twitter/login [get]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Generate random state
	state := "random-state" // TODO: Implement proper state management

	// Get authorization URL
	authURL := h.service.GetAuthURL(state)

	return c.Redirect(authURL)
}

// Callback handles Twitter OAuth2 callback
// @Summary Handle Twitter OAuth2 callback
// @Description Processes OAuth2 callback and returns JWT token
// @Tags auth
// @Produce json
// @Param code query string true "OAuth2 code"
// @Param state query string true "OAuth2 state"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} exception.ErrorResponse
// @Failure 401 {object} exception.ErrorResponse
// @Router /auth/twitter/callback [get]
func (h *AuthHandler) Callback(c *fiber.Ctx) error {
	// Get code from query
	code := c.Query("code")
	if code == "" {
		return exception.BadRequest("Missing code parameter")
	}

	// Verify state
	state := c.Query("state")
	if state == "" {
		return exception.BadRequest("Missing state parameter")
	}
	// TODO: Validate state

	// Exchange code for token
	token, err := h.service.ExchangeCode(context.Background(), code)
	if err != nil {
		return err
	}

	// Get user info
	user, err := h.service.GetUserInfo(context.Background(), token)
	if err != nil {
		return err
	}

	// Generate JWT token
	jwtToken, err := h.service.GenerateToken(user)
	if err != nil {
		return err
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    jwtToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.Status(http.StatusOK).JSON(LoginResponse{
		Token: jwtToken,
		User:  *user,
	})
}
