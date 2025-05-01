package config

import "time"

// AuthConfig holds authentication configuration values
type AuthConfig struct {
	JWTSecret        string
	TokenExpiration  time.Duration
	PasswordHashCost int
}

// NewAuthConfig creates a default authentication configuration
func NewAuthConfig() *AuthConfig {
	return &AuthConfig{
		JWTSecret:        "your-jwt-secret-change-me-in-production",
		TokenExpiration:  time.Hour * 24, // 24 hours
		PasswordHashCost: 10,
	}
}
