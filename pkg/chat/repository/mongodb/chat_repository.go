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

type ChatRepository struct {
	collection *mongo.Collection
}

func NewChatRepository(db *mongo.Database) (repository.ChatRepository, error) {
	repo := &ChatRepository{
		collection: db.Collection("chats"),
	}

	if err := repo.ensureIndexes(context.Background()); err != nil {
		return nil, err
	}

	return repo, nil
}

// ensureIndexes creates all necessary indexes for the chat collection
func (r *ChatRepository) ensureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "chatroom", Value: 1},
				{Key: "createdTimestamp", Value: -1},
			},
			Options: options.Index().SetName("chatroom_timestamp"),
		},
		{
			Keys: bson.D{
				{Key: "sender", Value: 1},
				{Key: "createdTimestamp", Value: -1},
			},
			Options: options.Index().SetName("sender_timestamp"),
		},
		{
			Keys: bson.D{
				{Key: "receiver", Value: 1},
				{Key: "createdTimestamp", Value: -1},
			},
			Options: options.Index().SetName("receiver_timestamp"),
		},
	}

	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	_, err := r.collection.Indexes().CreateMany(ctx, indexes, opts)
	return err
}

// Get retrieves a single chat message by ID
func (r *ChatRepository) Get(ctx context.Context, id string) (*entity.Chat, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var chat entity.Chat
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&chat)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &chat, nil
}

// GetPopulated retrieves a single chat message with populated user references
func (r *ChatRepository) GetPopulated(ctx context.Context, id string) (*entity.ChatPopulated, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": objectID}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "sender",
			"foreignField": "_id",
			"as":           "sender",
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "receiver",
			"foreignField": "_id",
			"as":           "receiver",
		}}},
		{{Key: "$addFields", Value: bson.M{
			"sender":   bson.M{"$arrayElemAt": []interface{}{"$sender", 0}},
			"receiver": bson.M{"$arrayElemAt": []interface{}{"$receiver", 0}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chats []entity.ChatPopulated
	if err = cursor.All(ctx, &chats); err != nil {
		return nil, err
	}

	if len(chats) == 0 {
		return nil, nil
	}

	return &chats[0], nil
}

// GetAll retrieves multiple chat messages with filtering and pagination
func (r *ChatRepository) GetAll(ctx context.Context, filter repository.ChatFilter, pag pagination.Pagination) ([]*entity.Chat, int64, error) {
	query := bson.M{}
	if filter.ChatroomID != "" {
		query["chatroom"] = filter.ChatroomID
	}
	if filter.SenderID != "" {
		query["sender"] = filter.SenderID
	}
	if filter.ReceiverID != "" {
		query["receiver"] = filter.ReceiverID
	}
	if filter.StartTime != nil {
		query["createdTimestamp"] = bson.M{"$gte": *filter.StartTime}
	}
	if filter.EndTime != nil {
		if _, exists := query["createdTimestamp"]; exists {
			query["createdTimestamp"].(bson.M)["$lte"] = *filter.EndTime
		} else {
			query["createdTimestamp"] = bson.M{"$lte": *filter.EndTime}
		}
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "createdTimestamp", Value: -1}}).
		SetSkip(int64((pag.Page - 1) * pag.Limit)).
		SetLimit(int64(pag.Limit))

	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var chats []*entity.Chat
	if err = cursor.All(ctx, &chats); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return chats, total, nil
}

// GetAllPopulated retrieves multiple chat messages with populated user references
func (r *ChatRepository) GetAllPopulated(ctx context.Context, filter repository.ChatFilter, pag pagination.Pagination) ([]*entity.ChatPopulated, int64, error) {
	matchStage := bson.M{}
	if filter.ChatroomID != "" {
		matchStage["chatroom"] = filter.ChatroomID
	}
	if filter.SenderID != "" {
		matchStage["sender"] = filter.SenderID
	}
	if filter.ReceiverID != "" {
		matchStage["receiver"] = filter.ReceiverID
	}
	if filter.StartTime != nil {
		matchStage["createdTimestamp"] = bson.M{"$gte": *filter.StartTime}
	}
	if filter.EndTime != nil {
		if _, exists := matchStage["createdTimestamp"]; exists {
			matchStage["createdTimestamp"].(bson.M)["$lte"] = *filter.EndTime
		} else {
			matchStage["createdTimestamp"] = bson.M{"$lte": *filter.EndTime}
		}
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchStage}},
		{{Key: "$sort", Value: bson.D{{Key: "createdTimestamp", Value: -1}}}},
		{{Key: "$skip", Value: int64((pag.Page - 1) * pag.Limit)}},
		{{Key: "$limit", Value: int64(pag.Limit)}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "sender",
			"foreignField": "_id",
			"as":           "sender",
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "receiver",
			"foreignField": "_id",
			"as":           "receiver",
		}}},
		{{Key: "$addFields", Value: bson.M{
			"sender":   bson.M{"$arrayElemAt": []interface{}{"$sender", 0}},
			"receiver": bson.M{"$arrayElemAt": []interface{}{"$receiver", 0}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var chats []*entity.ChatPopulated
	if err = cursor.All(ctx, &chats); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, matchStage)
	if err != nil {
		return nil, 0, err
	}

	return chats, total, nil
}

// Create stores a new chat message
func (r *ChatRepository) Create(ctx context.Context, chat *entity.Chat) error {
	if chat.ID == "" {
		chat.ID = primitive.NewObjectID().Hex()
	}

	_, err := r.collection.InsertOne(ctx, chat)
	return err
}

// Update modifies an existing chat message
func (r *ChatRepository) Update(ctx context.Context, chat *entity.Chat) error {
	objectID, err := primitive.ObjectIDFromHex(chat.ID)
	if err != nil {
		return err
	}

	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objectID}, chat)
	return err
}

// Delete removes a chat message
func (r *ChatRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
