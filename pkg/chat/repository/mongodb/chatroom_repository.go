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

type ChatroomRepository struct {
	collection *mongo.Collection
}

func NewChatroomRepository(db *mongo.Database) (repository.ChatroomRepository, error) {
	repo := &ChatroomRepository{
		collection: db.Collection("chatrooms"),
	}

	if err := repo.ensureIndexes(context.Background()); err != nil {
		return nil, err
	}

	return repo, nil
}

// ensureIndexes creates all necessary indexes for the chatroom collection
func (r *ChatroomRepository) ensureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "participants.user", Value: 1},
				{Key: "createdTimestamp", Value: -1},
			},
			Options: options.Index().SetName("participants_timestamp"),
		},
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
				{Key: "createdTimestamp", Value: -1},
			},
			Options: options.Index().SetName("type_timestamp"),
		},
	}

	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	_, err := r.collection.Indexes().CreateMany(ctx, indexes, opts)
	return err
}

// Get retrieves a single chatroom by ID
func (r *ChatroomRepository) Get(ctx context.Context, id string) (*entity.Chatroom, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var chatroom entity.Chatroom
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&chatroom)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &chatroom, nil
}

// GetPopulated retrieves a single chatroom with populated user references
func (r *ChatroomRepository) GetPopulated(ctx context.Context, id string) (*entity.ChatroomPopulated, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": objectID}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "lastSender",
			"foreignField": "_id",
			"as":           "lastSenderUser",
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "participants.user",
			"foreignField": "_id",
			"as":           "participantUsers",
		}}},
		{{Key: "$addFields", Value: bson.M{
			"lastSender": bson.M{"$arrayElemAt": []interface{}{"$lastSenderUser", 0}},
			"participants": bson.M{
				"$map": bson.M{
					"input": "$participants",
					"as":    "participant",
					"in": bson.M{
						"_id":                 "$$participant._id",
						"user":                bson.M{"$arrayElemAt": []interface{}{"$participantUsers", 0}},
						"role":                "$$participant.role",
						"joinedTimestamp":     "$$participant.joinedTimestamp",
						"mutedUntilTimestamp": "$$participant.mutedUntilTimestamp",
					},
				},
			},
		}}},
		{{Key: "$project", Value: bson.M{"lastSenderUser": 0, "participantUsers": 0}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chatrooms []entity.ChatroomPopulated
	if err = cursor.All(ctx, &chatrooms); err != nil {
		return nil, err
	}

	if len(chatrooms) == 0 {
		return nil, nil
	}

	return &chatrooms[0], nil
}

// GetAll retrieves multiple chatrooms with filtering and pagination
func (r *ChatroomRepository) GetAll(ctx context.Context, filter repository.ChatroomFilter, pag pagination.Pagination) ([]*entity.Chatroom, int64, error) {
	query := bson.M{}
	if filter.ParticipantID != "" {
		query["participants.user"] = filter.ParticipantID
	}
	if filter.Type != nil {
		query["type"] = *filter.Type
	}
	if filter.IsGroup != nil {
		query["isGroup"] = *filter.IsGroup
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

	var chatrooms []*entity.Chatroom
	if err = cursor.All(ctx, &chatrooms); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return chatrooms, total, nil
}

// GetAllPopulated retrieves multiple chatrooms with populated user references
func (r *ChatroomRepository) GetAllPopulated(ctx context.Context, filter repository.ChatroomFilter, pag pagination.Pagination) ([]*entity.ChatroomPopulated, int64, error) {
	matchStage := bson.M{}
	if filter.ParticipantID != "" {
		matchStage["participants.user"] = filter.ParticipantID
	}
	if filter.Type != nil {
		matchStage["type"] = *filter.Type
	}
	if filter.IsGroup != nil {
		matchStage["isGroup"] = *filter.IsGroup
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
			"localField":   "lastSender",
			"foreignField": "_id",
			"as":           "lastSenderUser",
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "participants.user",
			"foreignField": "_id",
			"as":           "participantUsers",
		}}},
		{{Key: "$addFields", Value: bson.M{
			"lastSender": bson.M{"$arrayElemAt": []interface{}{"$lastSenderUser", 0}},
			"participants": bson.M{
				"$map": bson.M{
					"input": "$participants",
					"as":    "participant",
					"in": bson.M{
						"_id":                 "$$participant._id",
						"user":                bson.M{"$arrayElemAt": []interface{}{"$participantUsers", 0}},
						"role":                "$$participant.role",
						"joinedTimestamp":     "$$participant.joinedTimestamp",
						"mutedUntilTimestamp": "$$participant.mutedUntilTimestamp",
					},
				},
			},
		}}},
		{{Key: "$project", Value: bson.M{"lastSenderUser": 0, "participantUsers": 0}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var chatrooms []*entity.ChatroomPopulated
	if err = cursor.All(ctx, &chatrooms); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, matchStage)
	if err != nil {
		return nil, 0, err
	}

	return chatrooms, total, nil
}

// Create stores a new chatroom
func (r *ChatroomRepository) Create(ctx context.Context, chatroom *entity.Chatroom) error {
	if chatroom.ID == "" {
		chatroom.ID = primitive.NewObjectID().Hex()
	}

	_, err := r.collection.InsertOne(ctx, chatroom)
	return err
}

// Update modifies an existing chatroom
func (r *ChatroomRepository) Update(ctx context.Context, chatroom *entity.Chatroom) error {
	objectID, err := primitive.ObjectIDFromHex(chatroom.ID)
	if err != nil {
		return err
	}

	_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objectID}, chatroom)
	return err
}

// Delete removes a chatroom
func (r *ChatroomRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// AddParticipant adds a participant to a chatroom
func (r *ChatroomRepository) AddParticipant(ctx context.Context, chatroomID string, participant entity.ChatroomParticipant) error {
	objectID, err := primitive.ObjectIDFromHex(chatroomID)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$push": bson.M{"participants": participant}},
	)
	return err
}

// RemoveParticipant removes a participant from a chatroom
func (r *ChatroomRepository) RemoveParticipant(ctx context.Context, chatroomID string, userID string) error {
	objectID, err := primitive.ObjectIDFromHex(chatroomID)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$pull": bson.M{"participants": bson.M{"user": userID}}},
	)
	return err
}

// UpdateParticipant updates a participant's properties in a chatroom
func (r *ChatroomRepository) UpdateParticipant(ctx context.Context, chatroomID string, participant entity.ChatroomParticipant) error {
	objectID, err := primitive.ObjectIDFromHex(chatroomID)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{
			"_id":               objectID,
			"participants.user": participant.User,
		},
		bson.M{"$set": bson.M{
			"participants.$.role":                participant.Role,
			"participants.$.mutedUntilTimestamp": participant.MutedUntilTimestamp,
		}},
	)
	return err
}
