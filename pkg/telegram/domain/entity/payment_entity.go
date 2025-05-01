package entity

import "time"

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	// PaymentStatusPending indicates a payment is pending
	PaymentStatusPending PaymentStatus = "pending"
	// PaymentStatusCompleted indicates a payment is completed
	PaymentStatusCompleted PaymentStatus = "completed"
	// PaymentStatusFailed indicates a payment has failed
	PaymentStatusFailed PaymentStatus = "failed"
	// PaymentStatusCancelled indicates a payment was cancelled
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// Payment represents a Telegram star payment
type Payment struct {
	ID            string        `json:"id"`
	UserID        int64         `json:"user_id"`
	ChatID        int64         `json:"chat_id"`
	Amount        float64       `json:"amount"`
	Currency      string        `json:"currency"`
	Description   string        `json:"description"`
	Status        PaymentStatus `json:"status"`
	InvoiceID     string        `json:"invoice_id,omitempty"`
	ProviderToken string        `json:"-"` // Provider token is sensitive and not stored in JSON
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	ClientID      string        `json:"client_id"`
}
