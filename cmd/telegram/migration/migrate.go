package main

import (
	"app/pkg/database/mongodb"
	"app/pkg/telegram/config"
	"app/pkg/telegram/domain/entity"
	repository "app/pkg/telegram/repository/mongodb"
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
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

	db := mongoClient.Database(cfg.MongoDB.Database)
	ctx := context.Background()

	// Run migrations
	if err := runMigrations(ctx, db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
}

func runMigrations(ctx context.Context, db *mongo.Database) error {
	log.Println("Starting migrations...")
	start := time.Now()

	// Create repositories
	userRepo, err := repository.NewUserRepository(db)
	if err != nil {
		return fmt.Errorf("failed to create user repository: %v", err)
	}

	clientRepo, err := repository.NewClientRepository(db)
	if err != nil {
		return fmt.Errorf("failed to create client repository: %v", err)
	}

	// Create sample clients
	now := time.Now().UnixMilli()
	clients := []*entity.Client{
		{
			Token:            "test_bot_token_1",
			Username:         "test_bot_1",
			Name:             "Test Bot 1",
			Description:      "Test bot for development purposes",
			BotType:          "support",
			WebhookURL:       "https://api.example.com/webhook/bot1",
			Status:           "active",
			MaxConnections:   40,
			AllowedUpdates:   []string{"message", "callback_query"},
			CreatedTimestamp: now,
			UpdatedTimestamp: now,
		},
		{
			Token:            "test_bot_token_2",
			Username:         "test_bot_2",
			Name:             "Test Bot 2",
			Description:      "Another test bot for development",
			BotType:          "notification",
			WebhookURL:       "https://api.example.com/webhook/bot2",
			Status:           "active",
			MaxConnections:   40,
			AllowedUpdates:   []string{"message", "callback_query", "channel_post"},
			CreatedTimestamp: now,
			UpdatedTimestamp: now,
		},
	}

	log.Println("Creating sample clients...")
	for _, client := range clients {
		if err := clientRepo.Create(ctx, client); err != nil {
			return fmt.Errorf("failed to create client: %v", err)
		}
	}

	// Create sample users
	users := []*entity.User{
		{
			ClientID:         "test_client",
			TelegramID:       123456789,
			Username:         "user1",
			FirstName:        "John",
			LastName:         "Doe",
			LanguageCode:     "en",
			IsBot:            false,
			IsPremium:        false,
			Status:           "active",
			CreatedTimestamp: now,
			UpdatedTimestamp: now,
		},
		{
			ClientID:         "test_client",
			TelegramID:       987654321,
			Username:         "user2",
			FirstName:        "Jane",
			LastName:         "Smith",
			LanguageCode:     "en",
			IsBot:            false,
			IsPremium:        false,
			Status:           "active",
			CreatedTimestamp: now,
			UpdatedTimestamp: now,
		},
	}

	log.Println("Creating sample users...")
	for _, user := range users {
		if err := userRepo.Create(ctx, user); err != nil {
			return fmt.Errorf("failed to create user: %v", err)
		}
	}

	log.Printf("Migration completed in %v\n", time.Since(start))
	return nil
}
