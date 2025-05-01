package config

import (
	"app/pkg/database"
	"app/pkg/fiber"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server  fiber.ServerConfig      `yaml:"server" env-prefix:"SERVER_"`
	MongoDB database.DatabaseConfig `yaml:"mongodb" env-prefix:"MONGODB_"`
	Redis   database.DatabaseConfig `yaml:"redis" env-prefix:"REDIS_"`
}

// LoadConfig loads configuration for a specific service
func LoadConfig[T any](configPath string) (*T, error) {
	var cfg T
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		// If config file not found, try to read from environment variables only
		if os.IsNotExist(err) {
			if err := cleanenv.ReadEnv(&cfg); err != nil {
				return nil, fmt.Errorf("error reading environment variables: %v", err)
			}
			return &cfg, nil
		}
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	return &cfg, nil
}
