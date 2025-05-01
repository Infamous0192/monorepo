package repository

import (
	"app/pkg/telegram/domain/entity"
	"app/pkg/types/pagination"
	"context"
	"time"
)

// PaymentFilter represents filtering options for payment queries
type PaymentFilter struct {
	ID        string
	UserID    *int64
	ChatID    *int64
	Status    entity.PaymentStatus
	StartTime *time.Time
	EndTime   *time.Time
	MinAmount *float64
	MaxAmount *float64
	Currency  string
}

// PaymentRepository defines the interface for payment data operations
type PaymentRepository interface {
	// Get retrieves a single payment by ID
	Get(ctx context.Context, id string) (*entity.Payment, error)

	// GetAll retrieves multiple payments with filtering and pagination
	GetAll(ctx context.Context, filter PaymentFilter, pagination pagination.Pagination) ([]*entity.Payment, int64, error)

	// GetByUserID retrieves all payments for a specific user with pagination
	GetByUserID(ctx context.Context, userID int64, pagination pagination.Pagination) ([]*entity.Payment, int64, error)

	// Create stores a new payment record
	Create(ctx context.Context, payment *entity.Payment) error

	// Update modifies an existing payment
	Update(ctx context.Context, payment *entity.Payment) error

	// UpdateStatus updates the status of a payment
	UpdateStatus(ctx context.Context, id string, status entity.PaymentStatus) error

	// Delete removes a payment
	Delete(ctx context.Context, id string) error
}
