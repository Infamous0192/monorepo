package dto

// CreateClientRequest represents the request body for creating a client
type CreateClientRequest struct {
	Token          string   `json:"token" validate:"required"`
	Username       string   `json:"username" validate:"required"`
	Name           string   `json:"name" validate:"required"`
	Description    string   `json:"description"`
	BotType        string   `json:"botType" validate:"required"`
	WebhookURL     string   `json:"webhookUrl"`
	Status         string   `json:"status" validate:"omitempty,oneof=active inactive"`
	MaxConnections int      `json:"maxConnections"`
	AllowedUpdates []string `json:"allowedUpdates"`
}

// UpdateClientRequest represents the request body for updating a client
type UpdateClientRequest struct {
	Token          string   `json:"token" validate:"required"`
	Username       string   `json:"username" validate:"required"`
	Name           string   `json:"name" validate:"required"`
	Description    string   `json:"description"`
	BotType        string   `json:"botType" validate:"required"`
	WebhookURL     string   `json:"webhookUrl"`
	Status         string   `json:"status" validate:"required,oneof=active inactive"`
	MaxConnections int      `json:"maxConnections"`
	AllowedUpdates []string `json:"allowedUpdates"`
}
