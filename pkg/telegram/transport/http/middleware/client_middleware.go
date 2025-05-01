package middleware

import (
	"app/pkg/exception"
	"app/pkg/telegram/service/client"

	"github.com/gofiber/fiber/v2"
)

// ClientMiddleware handles client authentication in HTTP requests
type ClientMiddleware struct {
	clientService client.ClientService
}

// NewClientMiddleware creates a new instance of ClientMiddleware
func NewClientMiddleware(clientService client.ClientService) *ClientMiddleware {
	return &ClientMiddleware{
		clientService: clientService,
	}
}

// ValidateKey middleware validates the X-Client-Key header for client access
func (m *ClientMiddleware) ValidateKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientKey := c.Get("X-Client-Key")
		if clientKey == "" {
			return exception.Http(401, "X-Client-Key header is required")
		}

		// Validate client key
		client, err := m.clientService.GetByToken(c.Context(), clientKey)
		if err != nil {
			return err
		}

		if client == nil {
			return exception.Http(401, "Invalid client key")
		}

		// Store client in context for later use
		c.Locals("client", client)

		return c.Next()
	}
}
