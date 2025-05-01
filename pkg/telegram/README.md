# Telegram Bot Service

A flexible and scalable Telegram bot service that supports multiple bot instances with different behaviors.

## Directory Structure

```
pkg/telegram/
├── domain/
│   ├── entity/          # Domain entities
│   │   ├── client.go    # Bot client configuration
│   │   └── user.go      # Telegram user entity
│   └── repository/      # Repository interfaces
├── handler/
│   ├── bot_handler.go   # HTTP endpoints for bot management
│   ├── webhook_handler.go # Webhook endpoints for Telegram updates
│   └── bot_handlers.go  # Bot-specific command handlers
├── repository/
│   └── mongodb/         # MongoDB implementations of repositories
├── service/
│   ├── bot/            # Bot service implementation
│   │   ├── types.go    # Service interfaces
│   │   ├── bot_service.go # Single bot implementation
│   │   └── bot_manager.go # Multiple bots manager
│   └── client/         # Client service for bot configurations
└── README.md           # This file

```

## Features

- Multiple bot support with different behaviors
- Webhook and long polling support
- Thread-safe bot management
- Command handling system
- Callback query support
- MongoDB for persistence
- Redis for caching

## Usage

### 1. Creating a Bot Client

```go
client := &entity.Client{
    Token:    "YOUR_BOT_TOKEN",
    BotType:  "support", // or "news", "default"
    Username: "MyBot",
    Name:     "My Support Bot",
}
```

### 2. Starting a Bot

```go
botManager := bot.NewBotManager()
err := botManager.StartBot(ctx, client)
```

### 3. Handling Commands

```go
// Get bot instance
botInstance, err := botManager.GetBot(client.ID)
if err != nil {
    return err
}

// Set up command handlers
handler.SetupBotByType(client.BotType, botInstance)
```

### 4. Adding Custom Bot Type

1. Create handler function:
```go
func setupCustomBot(bot bot.BotService) {
    bot.Command("start", func(ctx context.Context, m *telebot.Message) error {
        return bot.SendMessage(ctx, m.Chat.ID, "Welcome to Custom Bot!")
    })
    
    // Add more commands...
}
```

2. Register in handler registry:
```go
var BotHandlerRegistry = map[string]BotHandlerSetup{
    "default": setupDefaultBot,
    "custom":  setupCustomBot,
}
```

## API Endpoints

### Bot Management

- `POST /api/v1/bots/start`
  ```json
  {
      "clientId": "bot1"
  }
  ```

- `POST /api/v1/bots/stop/{clientId}`
- `POST /api/v1/bots/stop-all`
- `GET /api/v1/bots/status/{clientId}`
- `GET /api/v1/bots/status`

### Webhooks

- `POST /api/v1/webhook/{clientId}`
  - Receives updates from Telegram

## Configuration

The service uses environment variables for configuration:

```yaml
server:
  port: ${SERVER_PORT}
  # ... other server configs

mongodb:
  host: ${MONGODB_HOST}
  # ... other MongoDB configs

redis:
  host: ${REDIS_HOST}
  # ... other Redis configs
```

## Docker Support

Build and run with Docker:

```bash
# Build and start all services
docker-compose up -d

# Build just the telegram service
docker-compose build telegram-service

# View logs
docker-compose logs -f telegram-service
```
