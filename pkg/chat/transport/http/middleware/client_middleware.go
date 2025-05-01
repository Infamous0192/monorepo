// Package middleware provides HTTP middleware components
package middleware

import (
	"app/pkg/chat/service/client"

	"github.com/gofiber/fiber/v2"
)

// ClientMiddleware handles client validation in HTTP requests
type ClientMiddleware struct {
	clientService client.ClientService
}

// NewClientMiddleware creates a new instance of ClientMiddleware
func NewClientMiddleware(clientService client.ClientService) *ClientMiddleware {
	return &ClientMiddleware{
		clientService: clientService,
	}
}

// ValidateClientKey middleware validates the X-Client-Key header and stores the client in context
func (m *ClientMiddleware) ValidateKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientKey := c.Get("X-Client-Key")

		client, err := m.clientService.ValidateKey(c.Context(), clientKey)
		if err != nil {
			return err
		}

		c.Locals("client", client)
		return c.Next()
	}
}
