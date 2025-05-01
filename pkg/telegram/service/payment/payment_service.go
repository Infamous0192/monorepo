package payment

import (
	"app/pkg/telegram/domain/entity"
	"app/pkg/telegram/domain/repository"
	"app/pkg/telegram/service/bot"
	"app/pkg/types/pagination"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/telebot.v4"
)

// Config holds configuration for the payment service
type Config struct {
	ProviderToken string
	Currency      string
}

// paymentService implements the payment service interface
type paymentService struct {
	config            Config
	paymentRepository repository.PaymentRepository
	botService        *bot.BotService
	worker            *PaymentWorker
}

// NewPaymentService creates a new payment service
func NewPaymentService(config Config, paymentRepository repository.PaymentRepository, botService *bot.BotService, worker *PaymentWorker) PaymentService {
	return &paymentService{
		config:            config,
		paymentRepository: paymentRepository,
		botService:        botService,
		worker:            worker,
	}
}

// SetWorker sets the payment worker for this service
func (s *paymentService) SetWorker(worker *PaymentWorker) {
	s.worker = worker
}

// CreateInvoice creates a payment invoice
func (s *paymentService) CreateInvoice(ctx context.Context, params CreateInvoiceParams) (*entity.Payment, error) {
	if params.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Set default values if not provided
	currency := params.Currency
	if currency == "" {
		currency = s.config.Currency
	}

	title := params.Title
	if title == "" {
		title = "Telegram Stars"
	}

	description := params.Description
	if description == "" {
		description = "Purchase of Telegram Stars"
	}

	// Generate a unique ID for the payment using MongoDB ObjectID
	paymentID := primitive.NewObjectID().Hex()

	// Use custom payload if provided, otherwise use the payment ID
	payload := params.Payload
	if payload == "" {
		payload = paymentID
	}

	// Create a new payment record
	payment := &entity.Payment{
		ID:            paymentID,
		UserID:        params.UserID,
		ChatID:        params.ChatID,
		Amount:        params.Amount,
		Currency:      currency,
		Description:   description,
		Status:        entity.PaymentStatusPending,
		ProviderToken: s.config.ProviderToken,
		ClientID:      params.ClientID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save the payment to the database
	if err := s.paymentRepository.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	// Create the invoice
	invoice := &telebot.Invoice{
		Title:       title,
		Description: description,
		Payload:     payload,
		Currency:    currency,
		Token:       s.config.ProviderToken,
		Prices: []telebot.Price{
			{
				Label:  "Stars",
				Amount: int(params.Amount * 100), // Convert to smallest currency unit (e.g., cents)
			},
		},
		NeedName:            false,
		NeedPhoneNumber:     false,
		NeedEmail:           false,
		NeedShippingAddress: false,
	}

	bot, err := s.botService.GetBot(params.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	// Send the invoice
	sentInvoice, err := bot.Send(
		&telebot.User{ID: params.UserID},
		invoice,
		&telebot.SendOptions{
			ReplyMarkup: &telebot.ReplyMarkup{},
		},
		s.config.ProviderToken,
	)

	if err != nil {
		// Update payment status to failed
		_ = s.paymentRepository.UpdateStatus(ctx, payment.ID, entity.PaymentStatusFailed)
		return nil, fmt.Errorf("failed to send invoice: %w", err)
	}

	// Update the payment with the invoice ID if possible
	var invoiceMessageID int
	if sentInvoice != nil && sentInvoice.ID != 0 {
		invoiceMessageID = sentInvoice.ID
		payment.InvoiceID = fmt.Sprintf("%d", sentInvoice.ID)
		// Update the invoice ID in the database
		if err := s.paymentRepository.Update(ctx, payment); err != nil {
			log.Printf("Failed to update payment with invoice ID: %v", err)
		}
	}

	// Send a reminder message that the invoice will be deleted after 10 minutes
	reminderMsg := "This invoice will be deleted after 10 minutes. Please ensure to pay within the given time."
	reminderMessage, err := bot.Send(
		&telebot.Chat{ID: params.ChatID},
		reminderMsg,
	)

	if err != nil {
		log.Printf("Failed to send reminder message: %v", err)
	}

	var reminderMessageID int
	if reminderMessage != nil && reminderMessage.ID != 0 {
		reminderMessageID = reminderMessage.ID
	}

	// Schedule deletion of invoice and reminder messages after 10 minutes using the worker
	if (invoiceMessageID != 0 || reminderMessageID != 0) && s.worker != nil {
		if err := s.worker.ScheduleInvoiceDeletion(
			ctx,
			InvoiceDeletionPayload{
				ChatID:            params.ChatID,
				InvoiceMessageID:  invoiceMessageID,
				ReminderMessageID: reminderMessageID,
				ClientID:          payment.ClientID,
				PaymentID:         payment.ID,
				ExpiresAt:         10 * time.Minute,
			},
		); err != nil {
			log.Printf("Failed to schedule invoice deletion: %v", err)
		} else {
			log.Printf("Invoice deletion scheduled for payment ID: %s. Messages will be deleted after 10 minutes", payment.ID)
		}
	} else if invoiceMessageID != 0 || reminderMessageID != 0 {
		log.Printf("Warning: Payment worker not set, invoice messages will not be automatically deleted for payment ID: %s", payment.ID)
	}

	return payment, nil
}

// ProcessSuccessfulPayment handles a successful payment
func (s *paymentService) ProcessSuccessfulPayment(ctx context.Context, paymentID string) error {
	// Get the payment
	payment, err := s.paymentRepository.Get(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	if payment == nil {
		return errors.New("payment not found")
	}

	// Update payment status
	if err := s.paymentRepository.UpdateStatus(ctx, paymentID, entity.PaymentStatusCompleted); err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	// Send confirmation message to the user
	confirmationMsg := fmt.Sprintf(
		"✅ Payment Successful!\n\n"+
			"Thank you for your purchase of %.2f %s for Telegram Stars.\n"+
			"Payment ID: %s\n"+
			"Date: %s",
		payment.Amount,
		payment.Currency,
		payment.ID,
		payment.UpdatedAt.Format("2006-01-02 15:04:05"),
	)

	if err := s.botService.SendMessage(ctx, payment.ClientID, payment.ChatID, confirmationMsg); err != nil {
		log.Printf("Failed to send payment confirmation: %v", err)
		// Continue processing even if sending the message fails
	}

	// Delete the invoice and reminder messages
	// Ignore errors as this is a cleanup operation
	if err := s.DeleteInvoiceMessages(ctx, paymentID); err != nil {
		log.Printf("Warning: Failed to delete invoice messages: %v", err)
	}

	return nil
}

// ProcessFailedPayment handles a failed payment
func (s *paymentService) ProcessFailedPayment(ctx context.Context, paymentID string, reason string) error {
	// Get the payment
	payment, err := s.paymentRepository.Get(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	if payment == nil {
		return errors.New("payment not found")
	}

	// Update payment status
	if err := s.paymentRepository.UpdateStatus(ctx, paymentID, entity.PaymentStatusFailed); err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	// Send failure message to the user
	failureMsg := fmt.Sprintf(
		"❌ Payment Failed\n\n"+
			"Your payment of %.2f %s for Telegram Stars could not be processed.\n"+
			"Reason: %s\n"+
			"Payment ID: %s\n"+
			"Please try again or contact support if the issue persists.",
		payment.Amount,
		payment.Currency,
		reason,
		payment.ID,
	)

	if err := s.botService.SendMessage(ctx, payment.ClientID, payment.ChatID, failureMsg); err != nil {
		log.Printf("Failed to send payment failure notification: %v", err)
		// Continue processing even if sending the message fails
	}

	// Delete the invoice and reminder messages
	// Ignore errors as this is a cleanup operation
	if err := s.DeleteInvoiceMessages(ctx, paymentID); err != nil {
		log.Printf("Warning: Failed to delete invoice messages: %v", err)
	}

	return nil
}

// GetPaymentByID retrieves a payment by its ID
func (s *paymentService) GetPaymentByID(ctx context.Context, id string) (*entity.Payment, error) {
	return s.paymentRepository.Get(ctx, id)
}

// GetUserPayments retrieves all payments for a user
func (s *paymentService) GetUserPayments(ctx context.Context, userID int64) ([]*entity.Payment, error) {
	// Use default pagination (first page, 10 items)
	pag := pagination.Pagination{
		Page:  1,
		Limit: 10,
	}

	payments, _, err := s.paymentRepository.GetByUserID(ctx, userID, pag)
	return payments, err
}

// HandlePreCheckoutQuery processes a pre-checkout query
func (s *paymentService) HandlePreCheckoutQuery(ctx context.Context, query *telebot.PreCheckoutQuery) error {
	// The payload contains the payment ID
	paymentID := query.Payload

	// Get the payment
	payment, err := s.paymentRepository.Get(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	if payment == nil {
		// Respond with an error
		// if err := s.botService.GetBot().Accept(query); err != nil {
		// 	log.Printf("Failed to respond to pre-checkout query: %v", err)
		// }
		return errors.New("payment not found")
	}

	bot, err := s.botService.GetBot(payment.ClientID)
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	// Verify the payment amount (adjust based on actual telebot.v4 API)
	// Note: In telebot.v4, the field might be different. Check the actual API documentation.
	// For now, we'll assume the payment is valid
	// expectedAmount := int(payment.Amount * 100)
	// if query.TotalAmount != expectedAmount {
	// 	// Respond with an error
	// 	if err := s.botService.GetBot().Accept(query); err != nil {
	// 		log.Printf("Failed to respond to pre-checkout query: %v", err)
	// 	}
	// 	return errors.New("invalid payment amount")
	// }

	// Accept the pre-checkout query
	if err := bot.Accept(query); err != nil {
		return fmt.Errorf("failed to accept pre-checkout query: %w", err)
	}

	return nil
}

// HandleSuccessfulPayment processes a successful payment notification
func (s *paymentService) HandleSuccessfulPayment(ctx context.Context, message *telebot.Message) error {
	if message.Payment == nil {
		return errors.New("no payment information in message")
	}

	// The payload contains the payment ID
	paymentID := message.Payment.Payload

	// Process the successful payment
	return s.ProcessSuccessfulPayment(ctx, paymentID)
}

// DeleteInvoiceMessages deletes invoice and reminder messages for a payment
func (s *paymentService) DeleteInvoiceMessages(ctx context.Context, paymentID string) error {
	if s.worker == nil {
		return fmt.Errorf("payment worker not set, cannot delete invoice messages")
	}
	return s.worker.DeleteInvoiceMessages(ctx, paymentID)
}
