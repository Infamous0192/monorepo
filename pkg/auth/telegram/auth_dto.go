package telegram

// LoginRequest represents the request body for Telegram login
type LoginRequest struct {
	InitData string `json:"initData" validate:"required"`
}

// LoginResponse represents the response for successful login
type LoginResponse struct {
	Token string       `json:"token"`
	User  TelegramUser `json:"user"`
}

// ValidationResponse represents the response for successful validation
type ValidationResponse struct {
	Valid bool         `json:"valid"`
	User  TelegramUser `json:"user,omitempty"`
}
