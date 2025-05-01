package entity

// User represents a Telegram user with client reference
type User struct {
	ID               string `bson:"_id,omitempty" json:"id,omitempty"`
	ClientID         string `bson:"clientId" json:"clientId"`
	TelegramID       int64  `bson:"telegramId" json:"telegramId"`
	Username         string `bson:"username" json:"username"`
	FirstName        string `bson:"firstName" json:"firstName"`
	LastName         string `bson:"lastName" json:"lastName"`
	LanguageCode     string `bson:"languageCode" json:"languageCode"`
	IsBot            bool   `bson:"isBot" json:"isBot"`
	IsPremium        bool   `bson:"isPremium" json:"isPremium"`
	Status           string `bson:"status" json:"status"`
	CreatedTimestamp int64  `bson:"createdTimestamp" json:"createdTimestamp"`
	UpdatedTimestamp int64  `bson:"updatedTimestamp" json:"updatedTimestamp"`
}
