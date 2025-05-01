package entity

// Chat represents a chat message
type Chat struct {
	ID               string  `bson:"_id,omitempty" json:"id,omitempty"`
	Message          string  `bson:"message" json:"message"`
	Sender           string  `bson:"sender" json:"sender"`     // Reference to Users collection
	Receiver         *string `bson:"receiver" json:"receiver"` // Null for group chats
	Chatroom         string  `bson:"chatroom" json:"chatroom"` // Can be string ID or Chatroom object
	CreatedTimestamp int64   `bson:"createdTimestamp" json:"createdTimestamp"`
	Premium          *bool   `bson:"premium,omitempty" json:"premium,omitempty"`
}

// ChatPopulated represents a chat message with populated user references
type ChatPopulated struct {
	ID               string `bson:"_id,omitempty" json:"id,omitempty"`
	Message          string `bson:"message" json:"message"`
	Sender           User   `bson:"sender" json:"sender"`     // Populated User object
	Receiver         *User  `bson:"receiver" json:"receiver"` // Populated User object, null for group chats
	Chatroom         string `bson:"chatroom" json:"chatroom"` // Can be string ID or Chatroom object
	CreatedTimestamp int64  `bson:"createdTimestamp" json:"createdTimestamp"`
	Premium          *bool  `bson:"premium,omitempty" json:"premium,omitempty"`
}
