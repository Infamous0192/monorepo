package payment

import (
	"app/pkg/database/redis"
	"app/pkg/telegram/service/bot"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gopkg.in/telebot.v4"
)

const (
	// Redis channels
	channelInvoiceDeletion = "payment:invoice_deletion"
	channelPaymentExpiry   = "payment:expiry"

	// Redis key prefixes
	keyPrefixInvoiceDeletion = "payment:delete_invoice:%s"
	keyPrefixPaymentExpiry   = "payment:expiry:%s"

	// Default expiration times
	defaultInvoiceExpiration = 10 * time.Minute
)

// PaymentWorker handles background tasks related to payments
type PaymentWorker struct {
	redisClient *redis.Client
	botManager  *bot.BotService
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewPaymentWorker creates a new payment worker
func NewPaymentWorker(redisClient *redis.Client, botManager *bot.BotService) *PaymentWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &PaymentWorker{
		redisClient: redisClient,
		botManager:  botManager,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start begins listening for payment-related tasks
func (w *PaymentWorker) Start() error {
	log.Println("Starting payment worker...")

	// Subscribe to all payment-related channels
	channels := []string{
		channelInvoiceDeletion,
		channelPaymentExpiry,
	}

	pubsub, err := w.redisClient.Subscribe(w.ctx, channels...)
	if err != nil {
		return fmt.Errorf("failed to subscribe to payment channels: %w", err)
	}

	// Start a goroutine to handle messages
	go func() {
		defer pubsub.Close()

		// Listen for messages
		for {
			select {
			case <-w.ctx.Done():
				log.Println("Payment worker stopped")
				return
			default:
				msg, err := pubsub.ReceiveMessage(w.ctx)
				if err != nil {
					log.Printf("Error receiving message: %v", err)
					continue
				}

				// Process the message based on the channel
				switch msg.Channel {
				case channelInvoiceDeletion:
					w.processInvoiceDeletion(msg.Payload)
				case channelPaymentExpiry:
					w.processPaymentExpiry(msg.Payload)
				default:
					log.Printf("Unknown channel: %s", msg.Channel)
				}
			}
		}
	}()

	// Start a periodic task to check for expired items
	go w.runPeriodicTasks()

	return nil
}

// Stop stops the worker
func (w *PaymentWorker) Stop() {
	w.cancel()
}

// ScheduleInvoiceDeletion schedules an invoice for deletion after the specified duration
func (w *PaymentWorker) ScheduleInvoiceDeletion(ctx context.Context, payload InvoiceDeletionPayload) error {
	if payload.ExpiresAt == 0 {
		payload.ExpiresAt = defaultInvoiceExpiration
	}

	// Convert to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal deletion payload: %w", err)
	}

	// Create a unique key for this deletion job
	deletionKey := fmt.Sprintf(keyPrefixInvoiceDeletion, payload.PaymentID)

	// Store the deletion info in Redis with the specified expiration
	if err := w.redisClient.Set(ctx, deletionKey, string(payloadBytes), payload.ExpiresAt); err != nil {
		return fmt.Errorf("failed to store deletion job in Redis: %w", err)
	}

	// Publish a message to the deletion channel
	if err := w.redisClient.Publish(ctx, channelInvoiceDeletion, deletionKey); err != nil {
		return fmt.Errorf("failed to publish deletion job: %w", err)
	}

	log.Printf("Invoice deletion job scheduled with key: %s. Messages will be deleted after %v", deletionKey, payload.ExpiresAt)
	return nil
}

// DeleteInvoiceMessages deletes invoice and reminder messages for a payment
func (w *PaymentWorker) DeleteInvoiceMessages(ctx context.Context, paymentID string) error {
	// Get the deletion job key
	deletionKey := fmt.Sprintf(keyPrefixInvoiceDeletion, paymentID)

	return w.processInvoiceDeletion(deletionKey)
}

// processInvoiceDeletion handles a single invoice deletion job
func (w *PaymentWorker) processInvoiceDeletion(key string) error {
	log.Printf("Processing invoice deletion: %s", key)

	// Get the deletion job data from Redis
	data, err := w.redisClient.Get(w.ctx, key)
	if err != nil {
		log.Printf("Failed to get deletion job data: %v", err)
		return fmt.Errorf("failed to get deletion job data: %w", err)
	}

	var job InvoiceDeletionPayload
	if err := json.Unmarshal([]byte(data), &job); err != nil {
		log.Printf("Failed to unmarshal deletion job data: %v", err)
		return fmt.Errorf("failed to unmarshal deletion job data: %w", err)
	}

	bot, err := w.botManager.GetBot(job.PaymentID)
	if err != nil {
		log.Printf("Failed to delete invoice message: %v", err)
		return fmt.Errorf("failed to delete invoice message: %w", err)
	}

	// Delete the invoice message
	if job.InvoiceMessageID != 0 {
		if err := bot.Delete(&telebot.Message{
			ID:   job.InvoiceMessageID,
			Chat: &telebot.Chat{ID: job.ChatID},
		}); err != nil {
			log.Printf("Failed to delete invoice message: %v", err)
		} else {
			log.Printf("Deleted invoice message %d for chat %d", job.InvoiceMessageID, job.ChatID)
		}
	}

	// Delete the reminder message
	if job.ReminderMessageID != 0 {
		if err := bot.Delete(&telebot.Message{
			ID:   job.ReminderMessageID,
			Chat: &telebot.Chat{ID: job.ChatID},
		}); err != nil {
			log.Printf("Failed to delete reminder message: %v", err)
		} else {
			log.Printf("Deleted reminder message %d for chat %d", job.ReminderMessageID, job.ChatID)
		}
	}

	// Delete the job from Redis
	if err := w.redisClient.Del(w.ctx, key); err != nil {
		log.Printf("Failed to delete job from Redis: %v", err)
	}

	return nil
}

// processPaymentExpiry handles payment expiration
func (w *PaymentWorker) processPaymentExpiry(key string) error {
	log.Printf("Processing payment expiry: %s", key)

	// Get the payment expiry data from Redis
	data, err := w.redisClient.Get(w.ctx, key)
	if err != nil {
		log.Printf("Failed to get payment expiry data: %v", err)
		return fmt.Errorf("failed to get payment expiry data: %w", err)
	}

	// Parse the job data
	var job struct {
		PaymentID string `json:"paymentID"`
		ExpiresAt int64  `json:"expiresAt"`
	}

	if err := json.Unmarshal([]byte(data), &job); err != nil {
		log.Printf("Failed to unmarshal payment expiry data: %v", err)
		return fmt.Errorf("failed to unmarshal payment expiry data: %w", err)
	}

	// Mark the payment as expired/cancelled in the database
	// This is a placeholder for actual implementation
	log.Printf("Payment %s has expired", job.PaymentID)

	// Delete the job from Redis
	if err := w.redisClient.Del(w.ctx, key); err != nil {
		log.Printf("Failed to delete payment expiry job from Redis: %v", err)
	}

	return nil
}

// runPeriodicTasks runs periodic maintenance tasks
func (w *PaymentWorker) runPeriodicTasks() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			w.processExpiredInvoices()
			w.processExpiredPayments()
		}
	}
}

// processExpiredInvoices checks for and processes any expired invoices
func (w *PaymentWorker) processExpiredInvoices() {
	log.Println("Processing expired invoices...")

	// Get all keys matching the deletion job pattern
	keys, err := w.redisClient.ScanKeys(w.ctx, keyPrefixInvoiceDeletion+"*")
	if err != nil {
		log.Printf("Failed to scan for expired invoices: %v", err)
		return
	}

	for _, key := range keys {
		if err := w.processInvoiceDeletion(key); err != nil {
			log.Printf("Error processing expired invoice: %v", err)
		}
	}
}

// processExpiredPayments checks for and processes any expired payments
func (w *PaymentWorker) processExpiredPayments() {
	log.Println("Processing expired payments...")

	// Get all keys matching the payment expiry pattern
	keys, err := w.redisClient.ScanKeys(w.ctx, keyPrefixPaymentExpiry+"*")
	if err != nil {
		log.Printf("Failed to scan for expired payments: %v", err)
		return
	}

	for _, key := range keys {
		if err := w.processPaymentExpiry(key); err != nil {
			log.Printf("Error processing expired payment: %v", err)
		}
	}
}
