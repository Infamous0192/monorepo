package bot

// BotService defines the interface for Telegram bot operations
// type BotService interface {
// 	// StartBot initializes and starts a new bot instance
// 	StartBot(ctx context.Context, client *entity.Client) error

// 	// StopBot stops and removes a bot instance
// 	StopBot(ctx context.Context, client *entity.Client) error

// 	// StopAllBots stops all running bot instances
// 	StopAllBots(ctx context.Context) error

// 	// GetBot retrieves a bot instance by client ID
// 	GetBot(clientID string) (*telebot.Bot, error)

// 	// GetAllBots returns all running bot instances
// 	GetAllBots() map[string]*telebot.Bot

// 	// SendMessage sends a text message to a specified chat
// 	SendMessage(ctx context.Context, clientID string, chatID int64, text string) error

// 	// SetWebhook sets up a webhook for the bot
// 	SetWebhook(ctx context.Context, client *entity.Client) error

// 	// RemoveWebhook removes the webhook for the bot
// 	RemoveWebhook(ctx context.Context, clientID string) error
// }
