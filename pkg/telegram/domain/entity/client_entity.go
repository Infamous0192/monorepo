package entity

// Client represents a Telegram bot client configuration
type Client struct {
	ID               string   `bson:"_id,omitempty" json:"id,omitempty"`
	Token            string   `bson:"token" json:"token"`
	Username         string   `bson:"username" json:"username"`
	Name             string   `bson:"name" json:"name"`
	Description      string   `bson:"description" json:"description"`
	BotType          string   `bson:"botType" json:"botType"`
	WebhookURL       string   `bson:"webhookUrl" json:"webhookUrl"`
	Status           string   `bson:"status" json:"status"`
	MaxConnections   int      `bson:"maxConnections" json:"maxConnections"`
	AllowedUpdates   []string `bson:"allowedUpdates" json:"allowedUpdates"`
	CreatedTimestamp int64    `bson:"createdTimestamp" json:"createdTimestamp"`
	UpdatedTimestamp int64    `bson:"updatedTimestamp" json:"updatedTimestamp"`
}
