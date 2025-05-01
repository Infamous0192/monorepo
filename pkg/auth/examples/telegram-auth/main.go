package main

import (
	"app/pkg/auth/telegram"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Add middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Initialize auth service and handlers
	jwtSecret := os.Getenv("JWT_SECRET")
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	authService := telegram.NewAuthService(botToken, jwtSecret, nil)
	authHandler := telegram.NewAuthHandler(authService)
	authMiddleware := telegram.NewAuthMiddleware(authService)

	// Register auth routes
	authHandler.RegisterRoutes(app)

	// Example protected routes
	api := app.Group("/api")

	// Protected route requiring just authentication
	protected := api.Group("/protected", authMiddleware.RequireAuth())
	protected.Get("/", func(c *fiber.Ctx) error {
		// Get user claims from context
		claims := c.Locals("user").(*telegram.Claims)
		return c.JSON(fiber.Map{
			"message": "Protected route",
			"user_id": claims.UserID,
			"name":    claims.FirstName,
		})
	})

	// Protected route requiring specific role
	admin := api.Group("/admin", authMiddleware.RequireAuth(), authMiddleware.RequireRole("admin"))
	admin.Get("/", func(c *fiber.Ctx) error {
		claims := c.Locals("user").(*telegram.Claims)
		return c.JSON(fiber.Map{
			"message": "Admin route",
			"user_id": claims.UserID,
			"name":    claims.FirstName,
		})
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}
