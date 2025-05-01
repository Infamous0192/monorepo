package mongodb

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/chat/domain/repository"
	"app/pkg/types/pagination"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClientRepository struct {
	collection *mongo.Collection
}

func NewClientRepository(db *mongo.Database) (repository.ClientRepository, error) {
	repo := &ClientRepository{
		collection: db.Collection("clients"),
	}

	if err := repo.ensureIndexes(context.Background()); err != nil {
		return nil, err
	}

	return repo, nil
}

// ensureIndexes creates all necessary indexes for the client collection
func (r *ClientRepository) ensureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "clientKey", Value: 1},
			},
			Options: options.Index().SetName("clientKey").SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "name", Value: 1},
			},
			Options: options.Index().SetName("name"),
		},
	}

	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	_, err := r.collection.Indexes().CreateMany(ctx, indexes, opts)
	return err
}

// Get retrieves a single client by ID
func (r *ClientRepository) Get(ctx context.Context, id string) (*entity.Client, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var client entity.Client
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &client, nil
}

// GetByKey retrieves a single client by client key
func (r *ClientRepository) GetByKey(ctx context.Context, clientKey string) (*entity.Client, error) {
	var client entity.Client
	err := r.collection.FindOne(ctx, bson.M{"clientKey": clientKey}).Decode(&client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &client, nil
}

// GetAll retrieves multiple clients with pagination
func (r *ClientRepository) GetAll(ctx context.Context, pag pagination.Pagination) ([]*entity.Client, int64, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "name", Value: 1}}).
		SetSkip(int64((pag.Page - 1) * pag.Limit)).
		SetLimit(int64(pag.Limit))

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var clients []*entity.Client
	if err = cursor.All(ctx, &clients); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return clients, total, nil
}

// Create stores a new client
func (r *ClientRepository) Create(ctx context.Context, client *entity.Client) error {
	if client.ID == "" {
		client.ID = primitive.NewObjectID().Hex()
	}

	now := time.Now().UnixMilli()
	client.CreatedTimestamp = now
	client.UpdatedTimestamp = now

	_, err := r.collection.InsertOne(ctx, client)
	return err
}

// Update modifies an existing client
func (r *ClientRepository) Update(ctx context.Context, client *entity.Client) error {
	objectID, err := primitive.ObjectIDFromHex(client.ID)
	if err != nil {
		return err
	}

	client.UpdatedTimestamp = time.Now().UnixMilli()

	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objectID}, client)
	return err
}

// Delete removes a client
func (r *ClientRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
