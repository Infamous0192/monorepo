package config

import (
	"app/pkg/config"
	"app/pkg/database"
	"app/pkg/fiber"
)

// AppConfig holds application specific configuration
type AppConfig struct {
	APIKey string `yaml:"api_key" env:"API_KEY"`
}

// ChatConfig holds chat service specific configuration
type ChatConfig struct {
	Server  fiber.ServerConfig      `yaml:"server" env-prefix:"SERVER_"`
	MongoDB database.DatabaseConfig `yaml:"mongodb" env-prefix:"MONGODB_"`
	Redis   database.DatabaseConfig `yaml:"redis" env-prefix:"REDIS_"`
	App     AppConfig               `yaml:"app" env-prefix:"APP_"`
}

// Load loads chat service configuration
func Load() (*ChatConfig, error) {
	// Load configuration
	cfg, err := config.LoadConfig[ChatConfig]("chat")
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
