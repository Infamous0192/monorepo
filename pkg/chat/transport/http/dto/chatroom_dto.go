package dto

// CreateChatroomRequest represents the request body for creating a chatroom
type CreateChatroomRequest struct {
	Name         string   `json:"name" validate:"required"`
	Type         string   `json:"type" validate:"required,oneof=public private squad"`
	IsGroup      bool     `json:"isGroup"`
	Participants []string `json:"participants,omitempty" validate:"omitempty,min=1"`
}

// UpdateChatroomRequest represents the request body for updating a chatroom
type UpdateChatroomRequest struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"type" validate:"required,oneof=public private squad"`
}

// AddParticipantRequest represents the request body for adding a participant to a chatroom
type AddParticipantRequest struct {
	UserID string `json:"userId" validate:"required"`
	Role   string `json:"role" validate:"required,oneof=member admin super_admin"`
}

// UpdateParticipantRequest represents the request body for updating a participant's role
type UpdateParticipantRequest struct {
	Role string `json:"role" validate:"required,oneof=member admin super_admin"`
}

// MuteParticipantRequest represents the request body for muting a participant
type MuteParticipantRequest struct {
	Duration int64 `json:"duration" validate:"required,min=1"` // Duration in minutes
}
