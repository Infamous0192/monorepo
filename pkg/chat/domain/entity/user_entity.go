package entity

// User represents a user in the chat system
type User struct {
	ID       string `bson:"_id,omitempty" json:"id,omitempty"`
	UserID   string `bson:"userId" json:"userId"`
	Name     string `bson:"name" json:"name"`
	Username string `bson:"username" json:"username"`
	Picture  string `bson:"picture" json:"picture"`
	Level    int    `bson:"level" json:"level"`
}
