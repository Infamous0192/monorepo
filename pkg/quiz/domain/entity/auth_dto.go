package entity

import "time"

// RegisterDTO represents the data transfer object for user registration
type RegisterDTO struct {
	Name      string    `json:"name" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	BirthDate time.Time `json:"birthDate" validate:"required"`
}

// LoginDTO represents the data transfer object for user login
type LoginDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AuthResponseDTO represents the data transfer object for authentication responses
type AuthResponseDTO struct {
	Token string `json:"token"`
	Creds *User  `json:"creds"`
}

// RegisterResponseDTO represents the data transfer object for registration responses
type RegisterResponseDTO AuthResponseDTO
