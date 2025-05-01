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

// PaymentConfig holds payment specific configuration
type PaymentConfig struct {
	ProviderToken string `yaml:"provider_token" env:"PAYMENT_PROVIDER_TOKEN"`
	Currency      string `yaml:"currency" env:"PAYMENT_CURRENCY"`
}

// TelegramConfig holds telegram service specific configuration
type TelegramConfig struct {
	Server  fiber.ServerConfig      `yaml:"server" env-prefix:"SERVER_"`
	MongoDB database.DatabaseConfig `yaml:"mongodb" env-prefix:"MONGODB_"`
	Redis   database.DatabaseConfig `yaml:"redis" env-prefix:"REDIS_"`
	App     AppConfig               `yaml:"app" env-prefix:"APP_"`
	Payment PaymentConfig           `yaml:"payment" env-prefix:"PAYMENT_"`
}

// Load loads telegram service configuration
func Load(configPath string) (*TelegramConfig, error) {
	// Load configuration
	cfg, err := config.LoadConfig[TelegramConfig](configPath)
	if err != nil {
		return nil, err
	}

	// Set default values
	if cfg.Payment.Currency == "" {
		cfg.Payment.Currency = "USD"
	}

	return cfg, nil
}
