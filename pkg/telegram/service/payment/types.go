package payment

import (
	"app/pkg/telegram/domain/entity"
	"context"
	"time"

	"gopkg.in/telebot.v4"
)

// CreateInvoiceParams defines parameters for creating a payment invoice
type CreateInvoiceParams struct {
	UserID      int64   `json:"userId"`
	ChatID      int64   `json:"chatId"`
	Amount      float64 `json:"amount"`
	ClientID    string  `json:"clientId"`
	Currency    string  `json:"currency,omitempty"`    // Optional, will use default if empty
	Title       string  `json:"title,omitempty"`       // Optional, will use default if empty
	Description string  `json:"description,omitempty"` // Optional, will use default if empty
	Payload     string  `json:"payload,omitempty"`     // Optional custom payload
}

// PaymentService defines the interface for Telegram payment operations
type PaymentService interface {
	// CreateInvoice creates a payment invoice
	CreateInvoice(ctx context.Context, params CreateInvoiceParams) (*entity.Payment, error)

	// ProcessSuccessfulPayment handles a successful payment
	ProcessSuccessfulPayment(ctx context.Context, paymentID string) error

	// ProcessFailedPayment handles a failed payment
	ProcessFailedPayment(ctx context.Context, paymentID string, reason string) error

	// GetPaymentByID retrieves a payment by its ID
	GetPaymentByID(ctx context.Context, id string) (*entity.Payment, error)

	// GetUserPayments retrieves all payments for a user
	GetUserPayments(ctx context.Context, userID int64) ([]*entity.Payment, error)

	// HandlePreCheckoutQuery processes a pre-checkout query
	HandlePreCheckoutQuery(ctx context.Context, query *telebot.PreCheckoutQuery) error

	// HandleSuccessfulPayment processes a successful payment notification
	HandleSuccessfulPayment(ctx context.Context, message *telebot.Message) error

	// DeleteInvoiceMessages deletes invoice and reminder messages for a payment
	DeleteInvoiceMessages(ctx context.Context, paymentID string) error

	// SetWorker sets the payment worker for this service
	SetWorker(worker *PaymentWorker)
}

// InvoiceDeletionPayload defines the payload for invoice deletion
type InvoiceDeletionPayload struct {
	ChatID            int64         `json:"chatId"`
	InvoiceMessageID  int           `json:"invoiceMessageId"`
	ReminderMessageID int           `json:"reminderMessageId"`
	PaymentID         string        `json:"paymentId"`
	ClientID          string        `json:"clientId"`
	ExpiresAt         time.Duration `json:"expiresAt"`
}
