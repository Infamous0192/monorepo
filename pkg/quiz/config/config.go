package config

import (
	"github.com/gofiber/fiber/v2"
)

// QuizConfig represents the main configuration structure for the quiz service
type QuizConfig struct {
	App      AppConfig      `yaml:"app" env-prefix:"APP_"`
	Server   ServerConfig   `yaml:"server" env-prefix:"SERVER_"`
	Database DatabaseConfig `yaml:"database" env-prefix:"DB_"`
}

// AppConfig contains application-specific configuration
type AppConfig struct {
	Name        string `yaml:"name" env:"NAME" env-default:"quiz-service"`
	Environment string `yaml:"environment" env:"ENVIRONMENT" env-default:"development"`
	APIKey      string `yaml:"api_key" env:"API_KEY" env-default:"default-api-key"`
}

// ServerConfig contains server-specific configuration
type ServerConfig struct {
	Port            string       `yaml:"port" env:"PORT" env-default:"8080"`
	ReadTimeout     int          `yaml:"read_timeout" env:"READ_TIMEOUT" env-default:"10"`
	WriteTimeout    int          `yaml:"write_timeout" env:"WRITE_TIMEOUT" env-default:"10"`
	ShutdownTimeout int          `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" env-default:"10"`
	FiberConfig     fiber.Config `yaml:"-"`
}

// DatabaseConfig contains database connection details
type DatabaseConfig struct {
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"PORT" env-default:"5432"`
	User     string `yaml:"user" env:"USER" env-default:"postgres"`
	Password string `yaml:"password" env:"PASSWORD" env-default:"postgres"`
	Name     string `yaml:"name" env:"NAME" env-default:"quiz_db"`
	SSLMode  string `yaml:"ssl_mode" env:"SSL_MODE" env-default:"disable"`
}
