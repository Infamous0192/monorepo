package dto

// CreateClientRequest represents the request body for creating a client
type CreateClientRequest struct {
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description"`
	ClientKey    string `json:"clientKey" validate:"required"`
	AuthEndpoint string `json:"authEndpoint" validate:"required,url"`
}

// UpdateClientRequest represents the request body for updating a client
type UpdateClientRequest struct {
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description"`
	ClientKey    string `json:"clientKey" validate:"required"`
	AuthEndpoint string `json:"authEndpoint" validate:"required,url"`
	Status       string `json:"status" validate:"required,oneof=active inactive"`
}
