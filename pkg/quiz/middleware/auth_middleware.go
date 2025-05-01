package middleware

import (
	"app/pkg/exception"
	"app/pkg/quiz/services/auth"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware handles authentication for quiz API endpoints
type AuthMiddleware struct {
	authService auth.AuthService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(authService auth.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth middleware ensures the request has a valid JWT token
func (m *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// Try getting from cookie
			authHeader = c.Cookies("jwt")
			if authHeader == "" {
				return exception.Http(401, "Authentication required")
			}
		} else {
			// Remove Bearer prefix if present
			authHeader = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// Validate token
		claims, err := m.authService.ValidateToken(authHeader)
		if err != nil {
			return err
		}

		// Store claims in context for use in handlers
		c.Locals("user", claims)
		return c.Next()
	}
}

// RequireRole middleware ensures the user has the required role
func (m *AuthMiddleware) RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// First require authentication
		err := m.RequireAuth()(c)
		if err != nil {
			return err
		}

		// Get claims from context
		claims := c.Locals("user").(*auth.Claims)
		if claims.Role != role {
			return exception.Http(403, "Insufficient permissions")
		}

		return c.Next()
	}
}

// RequireAdmin middleware ensures the user is an admin
func (m *AuthMiddleware) RequireAdmin() fiber.Handler {
	return m.RequireRole("admin")
}
