package mongodb

import (
	"app/pkg/telegram/domain/entity"
	"app/pkg/telegram/domain/repository"
	"app/pkg/types/pagination"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PaymentRepository implements repository.PaymentRepository for MongoDB
type PaymentRepository struct {
	collection *mongo.Collection
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *mongo.Database) (repository.PaymentRepository, error) {
	repo := &PaymentRepository{
		collection: db.Collection("payments"),
	}

	if err := repo.ensureIndexes(context.Background()); err != nil {
		return nil, err
	}

	return repo, nil
}

// ensureIndexes creates all necessary indexes for the payment collection
func (r *PaymentRepository) ensureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().SetName("user_id_created_at"),
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().SetName("status_created_at"),
		},
	}

	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	_, err := r.collection.Indexes().CreateMany(ctx, indexes, opts)
	return err
}

// Get retrieves a single payment by ID
func (r *PaymentRepository) Get(ctx context.Context, id string) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&payment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &payment, nil
}

// GetAll retrieves multiple payments with filtering and pagination
func (r *PaymentRepository) GetAll(ctx context.Context, filter repository.PaymentFilter, pag pagination.Pagination) ([]*entity.Payment, int64, error) {
	query := bson.M{}

	if filter.ID != "" {
		query["_id"] = filter.ID
	}

	if filter.UserID != nil {
		query["user_id"] = *filter.UserID
	}

	if filter.ChatID != nil {
		query["chat_id"] = *filter.ChatID
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.StartTime != nil {
		query["created_at"] = bson.M{"$gte": *filter.StartTime}
	}

	if filter.EndTime != nil {
		if _, exists := query["created_at"]; exists {
			query["created_at"].(bson.M)["$lte"] = *filter.EndTime
		} else {
			query["created_at"] = bson.M{"$lte": *filter.EndTime}
		}
	}

	if filter.MinAmount != nil {
		query["amount"] = bson.M{"$gte": *filter.MinAmount}
	}

	if filter.MaxAmount != nil {
		if _, exists := query["amount"]; exists {
			query["amount"].(bson.M)["$lte"] = *filter.MaxAmount
		} else {
			query["amount"] = bson.M{"$lte": *filter.MaxAmount}
		}
	}

	if filter.Currency != "" {
		query["currency"] = filter.Currency
	}

	// Count total matching records
	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// Set up options for pagination and sorting
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64((pag.Page - 1) * pag.Limit)).
		SetLimit(int64(pag.Limit))

	// Execute the query
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Process the results
	var payments []*entity.Payment
	if err = cursor.All(ctx, &payments); err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// GetByUserID retrieves all payments for a specific user with pagination
func (r *PaymentRepository) GetByUserID(ctx context.Context, userID int64, pag pagination.Pagination) ([]*entity.Payment, int64, error) {
	// Create a filter with the user ID
	filter := repository.PaymentFilter{
		UserID: &userID,
	}

	// Use the GetAll method to handle filtering and pagination
	return r.GetAll(ctx, filter, pag)
}

// Create stores a new payment record
func (r *PaymentRepository) Create(ctx context.Context, payment *entity.Payment) error {
	if payment.ID == "" {
		payment.ID = primitive.NewObjectID().Hex()
	}

	now := time.Now()
	if payment.CreatedAt.IsZero() {
		payment.CreatedAt = now
	}
	payment.UpdatedAt = now

	// MongoDB document
	doc := bson.M{
		"_id":            payment.ID,
		"user_id":        payment.UserID,
		"chat_id":        payment.ChatID,
		"amount":         payment.Amount,
		"currency":       payment.Currency,
		"description":    payment.Description,
		"status":         payment.Status,
		"invoice_id":     payment.InvoiceID,
		"provider_token": payment.ProviderToken,
		"created_at":     payment.CreatedAt,
		"updated_at":     payment.UpdatedAt,
	}

	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

// Update modifies an existing payment
func (r *PaymentRepository) Update(ctx context.Context, payment *entity.Payment) error {
	payment.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"user_id":        payment.UserID,
			"chat_id":        payment.ChatID,
			"amount":         payment.Amount,
			"currency":       payment.Currency,
			"description":    payment.Description,
			"status":         payment.Status,
			"invoice_id":     payment.InvoiceID,
			"provider_token": payment.ProviderToken,
			"updated_at":     payment.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": payment.ID}, update)
	return err
}

// UpdateStatus updates the status of a payment
func (r *PaymentRepository) UpdateStatus(ctx context.Context, id string, status entity.PaymentStatus) error {
	now := time.Now()

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": now,
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete removes a payment
func (r *PaymentRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
