package middleware

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/chat/service/client"
	"app/pkg/exception"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	clientService client.ClientService
}

func NewAuthMiddleware(clientService client.ClientService) *AuthMiddleware {
	return &AuthMiddleware{
		clientService: clientService,
	}
}

// Authenticate middleware validates the Bearer token and stores the authenticated user in context
func (m *AuthMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get client from previous middleware
		client := c.Locals("client")
		if client == nil {
			return exception.InternalError("Client not found in context")
		}

		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return exception.Http(401, "Authorization header is required")
		}

		// Check bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return exception.Http(401, "Invalid authorization format. Use: Bearer <token>")
		}

		token := parts[1]
		if token == "" {
			return exception.Http(401, "Token is required")
		}

		// Get client ID from context
		clientData, ok := client.(*entity.Client)
		if !ok {
			return exception.InternalError("Invalid client data in context")
		}

		// Authenticate user
		user, err := m.clientService.Authenticate(c.Context(), clientData.ID, token)
		if err != nil {
			return err // Error is already wrapped with exception package
		}

		// Store authenticated user in context
		c.Locals("user", user)
		return c.Next()
	}
}
