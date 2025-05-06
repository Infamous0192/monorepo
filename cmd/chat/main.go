package main

import (
	_ "app/docs/api/chat" // Import generated docs
	"app/pkg/chat/config"
	repository "app/pkg/chat/repository/mongodb"
	"app/pkg/chat/service/chat"
	"app/pkg/chat/service/chatroom"
	"app/pkg/chat/service/client"
	"app/pkg/chat/service/user"
	"app/pkg/chat/transport/http/handler"
	"app/pkg/chat/transport/http/middleware"
	"app/pkg/chat/transport/ws"
	"app/pkg/database/mongodb"
	"app/pkg/database/redis"
	"app/pkg/fiber"
	sharedMiddleware "app/pkg/middleware"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	gofiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/ilyakaznacheev/cleanenv"
)

// @title Wonderverse Chat Service
// @version 1.0.0
// @description Chat Microservice for Wonderverse Apps
// @termsOfService http://swagger.io/terms/
// @contact.name Not Boring Company
// @contact.email infamous0192@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api
// @schemes http https
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Client-Key
// @description API Key for Client
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Authorization For JWT
func main() {
	var cfg config.ChatConfig
	var configPath = flag.String("config", filepath.Join("cmd", "chat", "config", "config.yml"), "path to config file")

	flag.Parse()

	// Load configuration
	if err := cleanenv.ReadConfig(*configPath, &cfg); err != nil {
		if os.IsNotExist(err) {
			if err := cleanenv.ReadEnv(&cfg); err != nil {
				log.Fatalf("error reading environment variables: %v", err)
			}
		}
		log.Fatalf("error reading config file: %v", err)
	}

	// Initialize MongoDB connection
	mongoClient := mongodb.NewClient(&cfg.MongoDB)
	if err := mongoClient.Connect(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer mongoClient.Disconnect()

	// Initialize Redis connection
	redisClient := redis.NewClient(cfg.Redis)
	if err := redisClient.Connect(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Disconnect()

	// Create repositories
	db := mongoClient.Database(cfg.MongoDB.Database)
	userRepo, err := repository.NewUserRepository(db)
	if err != nil {
		log.Fatalf("Failed to create user repository: %v", err)
	}
	clientRepo, err := repository.NewClientRepository(db)
	if err != nil {
		log.Fatalf("Failed to create client repository: %v", err)
	}
	chatRepo, err := repository.NewChatRepository(db)
	if err != nil {
		log.Fatalf("Failed to create chat repository: %v", err)
	}
	chatroomRepo, err := repository.NewChatroomRepository(db)
	if err != nil {
		log.Fatalf("Failed to create chatroom repository: %v", err)
	}

	// Create services
	userService := user.NewUserService(userRepo)
	clientService := client.NewClientService(clientRepo, userRepo, redisClient)
	chatroomService := chatroom.NewChatroomService(chatroomRepo)
	chatService := chat.NewChatService(chatRepo, chatroomService)

	// Create middleware
	clientMiddleware := middleware.NewClientMiddleware(clientService)
	authMiddleware := middleware.NewAuthMiddleware(clientService)
	adminMiddleware := sharedMiddleware.NewKeyMiddleware(cfg.App.APIKey)
	errorHandler := sharedMiddleware.NewErrorMiddleware()

	// Create WebSocket hub
	hub := ws.NewHub(redisClient, chatService, chatroomService)
	go hub.Run()

	// Create handlers
	userHandler := handler.NewUserHandler(userService, clientMiddleware, authMiddleware)
	clientHandler := handler.NewClientHandler(clientService, adminMiddleware)
	chatHandler := handler.NewChatHandler(chatService, clientMiddleware, authMiddleware)
	chatroomHandler := handler.NewChatroomHandler(chatroomService, clientMiddleware, authMiddleware)
	wsHandler := ws.NewHandler(redisClient, clientService, chatService, chatroomService)

	// API Custom error handler
	cfg.Server.ErrorHandler = errorHandler.Handler()

	// Create Fiber app with configuration
	app := fiber.NewServer(&cfg.Server)

	// Register routes
	api := app.Group("/api")
	userHandler.RegisterRoutes(api)
	clientHandler.RegisterRoutes(api)
	chatHandler.RegisterRoutes(api)
	chatroomHandler.RegisterRoutes(api)
	wsHandler.RegisterRoutes(api)

	// Swagger documentation route
	app.Get("/docs/*", swagger.HandlerDefault)

	app.Get("/health", func(c *gofiber.Ctx) error {
		return c.Status(gofiber.StatusOK).JSON(gofiber.Map{
			"status":      "ok",
			"time":        time.Now().Format(time.RFC3339),
			"service":     "chat",
			"version":     "1.0.2",
			"environment": "development",
			"uptime":      "just started",
			"hot_reload":  "working with absolute paths",
			"fixed":       true,
			"message":     "Hot reloading is working with absolute paths in Air config!",
		})
	})

	// Start server
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
			log.Printf("Server error: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
	}
}
