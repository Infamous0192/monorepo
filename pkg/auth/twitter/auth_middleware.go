package twitter

import (
	"app/pkg/exception"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	service *AuthService
}

// NewAuthMiddleware creates a new Twitter authentication middleware
func NewAuthMiddleware(service *AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		service: service,
	}
}

// RequireAuth middleware ensures the request has a valid JWT token
func (m *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		auth := c.Get("Authorization")
		if auth == "" {
			// Try getting from cookie
			auth = c.Cookies("jwt")
			if auth == "" {
				return exception.Http(401, "Missing authentication")
			}
		} else {
			// Remove Bearer prefix
			auth = strings.TrimPrefix(auth, "Bearer ")
		}

		// Validate token
		claims, err := m.service.ValidateToken(auth)
		if err != nil {
			return err
		}

		// Store claims in context
		c.Locals("user", claims)
		return c.Next()
	}
}

// RequireRole middleware ensures the user has the required role
func (m *AuthMiddleware) RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals("user").(*Claims)
		if claims.Role != role {
			return exception.Http(403, "Insufficient permissions")
		}
		return c.Next()
	}
}
