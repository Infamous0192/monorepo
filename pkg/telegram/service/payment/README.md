# Payment Service

This package implements a payment service for Telegram bots, handling invoice creation, payment processing, and automatic cleanup of expired invoices.

## Components

### PaymentService

The `PaymentService` interface defines methods for:
- Creating invoices
- Processing successful and failed payments
- Retrieving payment information
- Deleting invoice messages

### PaymentWorker

The `PaymentWorker` is a background service that:
- Listens for payment-related events via Redis PubSub
- Automatically deletes expired invoice messages
- Processes payment expiration events
- Runs periodic maintenance tasks

## How It Works

### Invoice Creation and Deletion Flow

1. When an invoice is created:
   - The invoice is sent to the user
   - A reminder message is sent indicating the invoice will expire
   - Both messages are scheduled for deletion after a specified time (default: 10 minutes)
   - The deletion job is stored in Redis with the appropriate expiration

2. The `PaymentWorker` handles deletion in two ways:
   - Reactively: When a deletion event is published to the Redis channel
   - Proactively: By periodically scanning for expired invoices

3. When a payment is successful or fails:
   - The invoice and reminder messages are deleted immediately
   - Appropriate confirmation/failure messages are sent to the user

### Redis Keys and Channels

- **Channels**:
  - `payment:invoice_deletion`: For invoice deletion events
  - `payment:expiry`: For payment expiration events

- **Key Prefixes**:
  - `payment:delete_invoice:{paymentID}`: Stores invoice deletion job data
  - `payment:expiry:{paymentID}`: Stores payment expiration job data

### Job Data Structure

Invoice deletion jobs contain:
```json
{
  "chatID": 123456789,
  "invoiceMessageID": 100,
  "reminderMessageID": 101,
  "paymentID": "abc123",
  "expiresAt": 1620000000
}
```

Payment expiry jobs contain:
```json
{
  "paymentID": "abc123",
  "expiresAt": 1620000000
}
```

## Usage

### Creating the Payment Worker

```go
// Create dependencies
redisClient := redis.NewClient(redisConfig)
botService := bot.NewBotService(botConfig)
paymentRepo := repository.NewPaymentRepository(db)

// Create and start the worker
worker := payment.NewPaymentWorker(redisClient, botService, paymentRepo)
if err := worker.Start(); err != nil {
    log.Fatalf("Failed to start payment worker: %v", err)
}

// The worker will run in the background
```

### Scheduling Invoice Deletion

```go
// After sending an invoice and reminder message
err := worker.ScheduleInvoiceDeletion(
    ctx,
    chatID,
    invoiceMessageID,
    reminderMessageID,
    paymentID,
    10*time.Minute,
)
```

### Manually Deleting Invoice Messages

```go
// When a payment is successful or fails
err := worker.DeleteInvoiceMessages(ctx, paymentID)
```

## Error Handling

The worker implements robust error handling:
- Failed message deletions are logged but don't stop processing
- Redis connection issues are reported but the worker continues to operate
- Periodic tasks ensure that no expired invoices are missed, even if real-time processing fails 