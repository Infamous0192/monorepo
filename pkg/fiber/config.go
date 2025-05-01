package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type CorsConfig struct {
	AllowOrigins     string `yaml:"allow_origins" env:"CORS_ALLOW_ORIGINS" env-default:"*"`
	AllowMethods     string `yaml:"allow_methods" env:"CORS_ALLOW_METHODS" env-default:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowHeaders     string `yaml:"allow_headers" env:"CORS_ALLOW_HEADERS" env-default:"Origin,Content-Type,Accept,Authorization"`
	AllowCredentials bool   `yaml:"allow_credentials" env:"CORS_ALLOW_CREDENTIALS" env-default:"true"`
	ExposeHeaders    string `yaml:"expose_headers" env:"CORS_EXPOSE_HEADERS" env-default:"Content-Length,Content-Type"`
	MaxAge           int    `yaml:"max_age" env:"CORS_MAX_AGE" env-default:"24h"`
}

type ServerConfig struct {
	// Server Configuration
	Name        string `yaml:"name" env:"NAME" env-default:"Not Boring Company"`
	Version     string `yaml:"version" env:"VERSION" env-default:"dev"`
	Env         string `yaml:"env" env:"ENV" env-default:"development"`
	Url         string `yaml:"url" env:"URL"`
	Host        string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port        string `yaml:"port" env:"PORT" env-default:"3001"`
	Path        string `yaml:"path" env:"PATH"`
	Debug       bool   `yaml:"debug" env:"DEBUG" env-default:"true"`
	PublicPath  string `yaml:"public_path" env:"PUBLIC_PATH"`
	UploadPath  string `yaml:"upload_path" env:"UPLOAD_PATH"`
	ExecPath    bool   `yaml:"exec_path" env:"EXEC_PATH" env-default:"false"`
	UploadLimit int    `yaml:"upload_limit" env:"UPLOAD_LIMIT" env-default:"8"`

	// Middleware Configuration
	Cors          *CorsConfig `yaml:"cors"`
	SwaggerConfig *swagger.Config
	ErrorHandler  fiber.ErrorHandler
}
