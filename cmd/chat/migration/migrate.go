package main

import (
	"app/pkg/chat/config"
	"app/pkg/chat/domain/entity"
	repository "app/pkg/chat/repository/mongodb"
	"app/pkg/database/mongodb"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"go.mongodb.org/mongo-driver/mongo"
)

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

	chatRepo, err := repository.NewChatRepository(db)
	if err != nil {
		return fmt.Errorf("failed to create chat repository: %v", err)
	}

	chatroomRepo, err := repository.NewChatroomRepository(db)
	if err != nil {
		return fmt.Errorf("failed to create chatroom repository: %v", err)
	}

	// Create sample users
	users := []*entity.User{
		{
			UserID:   "user1",
			Name:     "John Doe",
			Username: "johndoe",
			Picture:  "https://example.com/avatar1.jpg",
			Level:    1,
		},
		{
			UserID:   "user2",
			Name:     "Jane Smith",
			Username: "janesmith",
			Picture:  "https://example.com/avatar2.jpg",
			Level:    2,
		},
		{
			UserID:   "user3",
			Name:     "Admin User",
			Username: "admin",
			Picture:  "https://example.com/avatar3.jpg",
			Level:    10,
		},
	}

	log.Println("Creating sample users...")
	for _, user := range users {
		if err := userRepo.Create(ctx, user); err != nil {
			return fmt.Errorf("failed to create user: %v", err)
		}
	}

	// Create sample client
	client := &entity.Client{
		Name:         "Test Client",
		Description:  "Test client for development",
		ClientKey:    "test_client_key",
		Status:       "active",
		AuthEndpoint: "http://localhost:8080/auth",
	}

	log.Println("Creating sample client...")
	if err := clientRepo.Create(ctx, client); err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	// Create sample chatroom
	now := time.Now().UnixMilli()
	chatroom := &entity.Chatroom{
		Name:             "General Chat",
		IsGroup:          true,
		Type:             entity.ChatroomTypePublic,
		CreatedTimestamp: now,
		MessagesCount:    0,
		Participants: []entity.ChatroomParticipant{
			{
				User:            users[0].ID, // John Doe
				Role:            entity.ParticipantRoleMember,
				JoinedTimestamp: now,
			},
			{
				User:            users[1].ID, // Jane Smith
				Role:            entity.ParticipantRoleMember,
				JoinedTimestamp: now,
			},
			{
				User:            users[2].ID, // Admin
				Role:            entity.ParticipantRoleSuperAdmin,
				JoinedTimestamp: now,
			},
		},
	}

	log.Println("Creating sample chatroom...")
	if err := chatroomRepo.Create(ctx, chatroom); err != nil {
		return fmt.Errorf("failed to create chatroom: %v", err)
	}

	// Create sample chat messages
	messages := []string{
		"Hello everyone!",
		"Hi! How are you all doing?",
		"Welcome to the general chat!",
	}

	log.Println("Creating sample chat messages...")
	for i, msg := range messages {
		chat := &entity.Chat{
			Message:          msg,
			Sender:           users[i].ID,
			Chatroom:         chatroom.ID,
			CreatedTimestamp: now + int64(i*1000), // Messages 1 second apart
		}
		if err := chatRepo.Create(ctx, chat); err != nil {
			return fmt.Errorf("failed to create chat message: %v", err)
		}
	}

	// Update chatroom with last message
	chatroom.LastMessage = &messages[len(messages)-1]
	chatroom.LastSender = &users[len(messages)-1].ID
	chatroom.LastMessageTimestamp = &now
	chatroom.MessagesCount = len(messages)

	if err := chatroomRepo.Update(ctx, chatroom); err != nil {
		return fmt.Errorf("failed to update chatroom: %v", err)
	}

	log.Printf("Migration completed in %v\n", time.Since(start))
	return nil
}
