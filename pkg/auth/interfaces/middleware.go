package interfaces

import "github.com/gofiber/fiber/v2"

type AuthMiddleware interface {
	RequireAuth() fiber.Handler
	RequireRole(roles ...string) fiber.Handler
	RequireAdmin() fiber.Handler
}
