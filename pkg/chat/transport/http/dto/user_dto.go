package dto

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	UserID   string `json:"userId" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Picture  string `json:"picture"`
	Level    int    `json:"level" validate:"min=0"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Picture  string `json:"picture"`
	Level    int    `json:"level" validate:"min=0"`
}
