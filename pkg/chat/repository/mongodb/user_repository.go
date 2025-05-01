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

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) (repository.UserRepository, error) {
	repo := &UserRepository{
		collection: db.Collection("users"),
	}

	if err := repo.ensureIndexes(context.Background()); err != nil {
		return nil, err
	}

	return repo, nil
}

// ensureIndexes creates all necessary indexes for the user collection
func (r *UserRepository) ensureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "userId", Value: 1},
			},
			Options: options.Index().SetName("userId").SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "username", Value: 1},
			},
			Options: options.Index().SetName("username").SetUnique(true),
		},
	}

	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	_, err := r.collection.Indexes().CreateMany(ctx, indexes, opts)
	return err
}

// Get retrieves a single user by ID
func (r *UserRepository) Get(ctx context.Context, id string) (*entity.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user entity.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// GetAll retrieves multiple users with pagination
func (r *UserRepository) GetAll(ctx context.Context, pag pagination.Pagination) ([]*entity.User, int64, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "username", Value: 1}}).
		SetSkip(int64((pag.Page - 1) * pag.Limit)).
		SetLimit(int64(pag.Limit))

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []*entity.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Create stores a new user
func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	if user.ID == "" {
		user.ID = primitive.NewObjectID().Hex()
	}

	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// Update modifies an existing user
func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return err
	}

	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objectID}, user)
	return err
}

// Delete removes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
