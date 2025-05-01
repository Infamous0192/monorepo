package entity

// Client represents an application client in the system
type Client struct {
	ID               string `bson:"_id,omitempty" json:"id,omitempty"`
	Name             string `bson:"name" json:"name"`
	Description      string `bson:"description" json:"description"`
	ClientKey        string `bson:"clientKey" json:"clientKey"`
	Status           string `bson:"status" json:"status"`
	AuthEndpoint     string `bson:"authEndpoint" json:"authEndpoint"`
	CreatedTimestamp int64  `bson:"createdTimestamp" json:"createdTimestamp"`
	UpdatedTimestamp int64  `bson:"updatedTimestamp" json:"updatedTimestamp"`
}
