package bot

import (
	"app/pkg/telegram/domain/entity"
	"context"
	"fmt"
	"sync"
	"time"

	"gopkg.in/telebot.v4"
)

// BotService handles multiple bot instances
type BotService struct {
	bots  map[string]*telebot.Bot
	mutex sync.RWMutex
}

// NewBotService creates a new instance of BotService
func NewBotService() *BotService {
	return &BotService{}
}

// StartBot initializes and starts a new bot instance
func (m *BotService) StartBot(ctx context.Context, client *entity.Client) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if bot already exists
	if _, exists := m.bots[client.ID]; exists {
		return fmt.Errorf("bot with client ID %s already running", client.ID)
	}

	settings := telebot.Settings{
		Token:  client.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	if client.WebhookURL != "" {
		webhook := &telebot.Webhook{
			Listen:         ":8443",
			MaxConnections: client.MaxConnections,
			Endpoint: &telebot.WebhookEndpoint{
				PublicURL: client.WebhookURL,
			},
		}
		settings.Poller = webhook
	}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		return fmt.Errorf("failed to create bot: %w", err)
	}

	go bot.Start()

	// Store the bot instance
	m.bots[client.ID] = bot

	return nil
}

// StopBot stops and removes a bot instance
func (m *BotService) StopBot(ctx context.Context, client *entity.Client) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	bot, exists := m.bots[client.ID]
	if !exists {
		return fmt.Errorf("bot with client ID %s not found", client.ID)
	}

	if bot != nil {
		bot.Stop()
		delete(m.bots, client.ID)
	}

	return nil
}

// StopAllBots stops all running bot instances
func (m *BotService) StopAllBots(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, bot := range m.bots {
		bot.Stop()
	}

	// Clear the bots map regardless of errors
	m.bots = make(map[string]*telebot.Bot)

	return nil
}

// GetBot retrieves a bot instance by client ID
func (m *BotService) GetBot(clientID string) (*telebot.Bot, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	bot, exists := m.bots[clientID]
	if !exists {
		return nil, fmt.Errorf("bot with client ID %s not found", clientID)
	}

	return bot, nil
}

// GetAllBots returns all running bot instances
func (m *BotService) GetAllBots() map[string]*telebot.Bot {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Create a copy to avoid external modifications
	bots := make(map[string]*telebot.Bot, len(m.bots))
	for k, v := range m.bots {
		bots[k] = v
	}

	return bots
}

// SendMessage sends a text message to a specified chat
func (s *BotService) SendMessage(ctx context.Context, clientID string, chatID int64, text string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	bot, err := s.GetBot(clientID)
	if err != nil {
		return fmt.Errorf("failed to get bot: %w", err)
	}

	chat := &telebot.Chat{ID: chatID}

	if _, err := bot.Send(chat, text); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// SetWebhook sets up a webhook for the bot
func (s *BotService) SetWebhook(ctx context.Context, client *entity.Client) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	bot, err := s.GetBot(client.ID)
	if err != nil {
		return fmt.Errorf("failed to get bot: %w", err)
	}

	webhook := &telebot.Webhook{
		Endpoint: &telebot.WebhookEndpoint{
			PublicURL: client.WebhookURL,
		},
	}

	if err := bot.SetWebhook(webhook); err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	return nil
}

// RemoveWebhook removes the webhook for the bot
func (s *BotService) RemoveWebhook(ctx context.Context, clientID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	bot, err := s.GetBot(clientID)
	if err != nil {
		return fmt.Errorf("failed to get bot: %w", err)
	}

	if err := bot.RemoveWebhook(); err != nil {
		return fmt.Errorf("failed to remove webhook: %w", err)
	}

	return nil
}
