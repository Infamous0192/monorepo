package telegram

import (
	"app/pkg/exception"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware handles Telegram authentication middleware
type AuthMiddleware struct {
	authService *AuthService
}

// NewAuthMiddleware creates a new instance of AuthMiddleware
func NewAuthMiddleware(authService *AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth middleware validates the JWT token in the Authorization header
func (m *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return exception.Http(401, "Authorization header is required")
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return exception.Http(401, "Invalid authorization format")
		}

		claims, err := m.authService.ValidateToken(parts[1])
		if err != nil {
			return err
		}

		// Store user claims in context
		c.Locals("user", claims)
		return c.Next()
	}
}

// RequireRole middleware ensures the user has required role/permission
func (m *AuthMiddleware) RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("user").(*Claims)
		if !ok {
			return exception.Http(401, "User not authenticated")
		}

		// Here you can implement your role checking logic
		// For example, check if user.ID is in admin list, etc.
		// For now, we'll just check if the user exists
		if claims == nil {
			return exception.Http(403, "Access denied")
		}

		return c.Next()
	}
}
