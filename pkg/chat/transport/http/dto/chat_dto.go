package dto

// SendMessageRequest represents the request body for sending a message to a chatroom
type SendMessageRequest struct {
	Message string `json:"message" validate:"required"`
}

// SendDirectMessageRequest represents the request body for sending a direct message
type SendDirectMessageRequest struct {
	ReceiverID string `json:"receiverId" validate:"required"`
	Message    string `json:"message" validate:"required"`
}

// UpdateChatRequest represents the request body for updating a chat message
type UpdateChatRequest struct {
	Message string `json:"message" validate:"required"`
}
