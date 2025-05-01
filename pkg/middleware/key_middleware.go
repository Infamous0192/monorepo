package middleware

import (
	"app/pkg/exception"
	"crypto/subtle"

	"github.com/gofiber/fiber/v2"
)

// KeyMiddleware handles API Key authentication in HTTP requests
type KeyMiddleware struct {
	apiKey string
}

// NewKeyMiddleware creates a new instance of KeyMiddleware
func NewKeyMiddleware(apiKey string) *KeyMiddleware {

	return &KeyMiddleware{
		apiKey: apiKey,
	}
}

// ValidateAPIKey middleware validates the X-API-Key header for key access
func (m *KeyMiddleware) ValidateKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return exception.Http(401, "X-API-Key header is required")
		}

		if subtle.ConstantTimeCompare([]byte(apiKey), []byte(m.apiKey)) != 1 {
			return exception.Http(401, "Invalid API Key")
		}

		return c.Next()
	}
}
