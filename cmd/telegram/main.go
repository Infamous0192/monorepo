package main

import (
	"app/pkg/database/mongodb"
	"app/pkg/database/redis"
	"app/pkg/fiber"
	sharedMiddleware "app/pkg/middleware"
	"app/pkg/telegram/config"
	repository "app/pkg/telegram/repository/mongodb"
	"app/pkg/telegram/service/bot"
	"app/pkg/telegram/service/client"
	"app/pkg/telegram/service/payment"
	httpHandler "app/pkg/telegram/transport/http/handler"
	httpMiddleware "app/pkg/telegram/transport/http/middleware"
	"app/pkg/types/pagination"
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
)

func main() {
	// Parse command line flags
	var configPath = flag.String("config", filepath.Join("cmd", "telegram", "config", "config.yml"), "path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
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
	clientRepo, err := repository.NewClientRepository(db)
	if err != nil {
		log.Fatalf("Failed to create client repository: %v", err)
	}

	// Create services
	clientService := client.NewClientService(clientRepo)
	botService := bot.NewBotService()

	// Create and start the payment worker
	var paymentWorker *payment.PaymentWorker
	ctx := context.Background()

	// Create middleware
	keyMiddleware := sharedMiddleware.NewKeyMiddleware(cfg.App.APIKey)
	clientMiddleware := httpMiddleware.NewClientMiddleware(clientService)
	errorHandler := sharedMiddleware.NewErrorMiddleware()

	// Create handlers
	botHandler := httpHandler.NewBotHandler(botService, clientService)
	webhookHandler := httpHandler.NewWebhookHandler(botService, clientService)
	clientHandler := httpHandler.NewClientHandler(clientService, keyMiddleware, clientMiddleware)

	// Create worker
	paymentWorker = payment.NewPaymentWorker(redisClient, botService)

	cfg.Server.ErrorHandler = errorHandler.Handler()

	// Create and configure Fiber app
	app := fiber.NewServer(&cfg.Server)

	api := app.Group("/api")

	// Register routes
	botHandler.RegisterRoutes(api)
	webhookHandler.RegisterRoutes(api)
	clientHandler.RegisterRoutes(api)

	app.Get("/health", func(c *gofiber.Ctx) error {
		fmt.Println(cfg.App)

		return c.Status(gofiber.StatusOK).JSON(gofiber.Map{
			"data":        "ok",
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

	// paymentService = payment.NewPaymentService(payment.Config{
	// 	ProviderToken: "",
	// 	Currency:      "XTR", // currency for telegram sstar
	// }, paymentRepo, botService, paymentWorker)

	// Start active bots
	activeClients, _, err := clientService.GetAll(ctx, pagination.Pagination{Limit: 100})
	if err != nil {
		log.Printf("Warning: Failed to get active clients: %v", err)
	} else {
		for _, client := range activeClients {
			// Start the bot
			if err := botService.StartBot(ctx, client); err != nil {
				log.Printf("Warning: Failed to start bot for client %s: %v", client.ID, err)
				continue
			}

			// Get the bot instance
			// botInstance, err := botService.GetBot(client.ID)
			// if err != nil {
			// 	log.Printf("Warning: Failed to get bot instance for client %s: %v", client.ID, err)
			// 	continue
			// }

			// Set up command handlers based on bot type
			botType := client.BotType
			if botType == "" {
				botType = "default" // Use default handlers if no type specified
			}

			// if err := botHandlers.SetupBotByType(botType, botInstance); err != nil {
			// 	log.Printf("Warning: Failed to set up handlers for bot %s: %v", client.ID, err)
			// 	continue
			// }

			log.Printf("Started bot type '%s' for client %s", botType, client.ID)

			// If this is a payment bot, initialize the payment worker for it
			// if botType == "payment" && paymentWorker == nil {
			// 	// Create payment service with this bot
			// 	paymentConfig := payment.Config{
			// 		ProviderToken: cfg.Payment.ProviderToken,
			// 		Currency:      cfg.Payment.Currency,
			// 	}
			// 	paymentService = payment.NewPaymentService(paymentConfig, botInstance, paymentRepo, redisClient)

			// 	// Create and start the worker
			// 	paymentWorker = payment.NewPaymentWorker(redisClient, botInstance, paymentRepo)
			// 	if err := paymentWorker.Start(); err != nil {
			// 		log.Printf("Warning: Failed to start payment worker: %v", err)
			// 	} else {
			// 		log.Println("Payment worker started successfully")

			// 		// Set the worker in the payment service
			// 		paymentService.SetWorker(paymentWorker)
			// 	}
			// }
		}
	}

	// If no payment bots were found but we still want to run the worker
	// for background cleanup of expired invoices
	// if paymentWorker == nil && len(activeClients) > 0 {
	// 	// Use the first available bot for the worker
	// 	firstBot, err := botManager.GetBot(activeClients[0].ID)
	// 	if err == nil {
	// 		// Create payment service with this bot
	// 		paymentConfig := payment.Config{
	// 			ProviderToken: cfg.Payment.ProviderToken,
	// 			Currency:      cfg.Payment.Currency,
	// 		}
	// 		paymentService = payment.NewPaymentService(paymentConfig, firstBot, paymentRepo, redisClient)

	// 		// Create and start the worker
	// 		paymentWorker = payment.NewPaymentWorker(redisClient, firstBot, paymentRepo)
	// 		if err := paymentWorker.Start(); err != nil {
	// 			log.Printf("Warning: Failed to start payment worker with fallback bot: %v", err)
	// 		} else {
	// 			log.Println("Payment worker started with fallback bot")

	// 			// Set the worker in the payment service
	// 			paymentService.SetWorker(paymentWorker)
	// 		}
	// 	}
	// }

	// Start server in a goroutine
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Stop the payment worker if it was started
	if paymentWorker != nil {
		paymentWorker.Stop()
		log.Println("Payment worker stopped")
	}

	// Stop all bots
	if err := botService.StopAllBots(ctx); err != nil {
		log.Printf("Error stopping bots: %v", err)
	}

	// Shutdown server with timeout
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Shutdown(); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	log.Println("Server stopped")
}
