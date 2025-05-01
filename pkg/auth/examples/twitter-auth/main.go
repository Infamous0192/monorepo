package main

import (
	"app/pkg/auth/twitter"
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

	// Get configuration from environment
	clientID := os.Getenv("TWITTER_CLIENT_ID")
	clientSecret := os.Getenv("TWITTER_CLIENT_SECRET")
	redirectURL := os.Getenv("TWITTER_REDIRECT_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	// Use default values for example (DO NOT use in production)
	if clientID == "" {
		clientID = "your-client-id"
	}
	if clientSecret == "" {
		clientSecret = "your-client-secret"
	}
	if redirectURL == "" {
		redirectURL = "http://localhost:3000/auth/twitter/callback"
	}
	if jwtSecret == "" {
		jwtSecret = "your-jwt-secret"
	}

	// Initialize auth components
	authService := twitter.NewAuthService(clientID, clientSecret, redirectURL, jwtSecret, nil)
	authHandler := twitter.NewAuthHandler(authService)
	authMiddleware := twitter.NewAuthMiddleware(authService)

	// Register auth routes
	authHandler.RegisterRoutes(app)

	// Example protected routes
	api := app.Group("/api")

	// Protected route requiring just authentication
	protected := api.Group("/protected", authMiddleware.RequireAuth())
	protected.Get("/", func(c *fiber.Ctx) error {
		// Get user claims from context
		claims := c.Locals("user").(*twitter.Claims)
		return c.JSON(fiber.Map{
			"message": "Protected route",
			"user_id": claims.UserID,
			"name":    claims.DisplayName,
		})
	})

	// Protected route requiring admin role
	admin := api.Group("/admin", authMiddleware.RequireAuth(), authMiddleware.RequireRole("admin"))
	admin.Get("/", func(c *fiber.Ctx) error {
		claims := c.Locals("user").(*twitter.Claims)
		return c.JSON(fiber.Map{
			"message": "Admin route",
			"user_id": claims.UserID,
			"name":    claims.DisplayName,
		})
	})

	// Start server
	log.Fatal(app.Listen(":3000"))
}
