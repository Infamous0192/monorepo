// Package docs contains Swagger/OpenAPI documentation for the API
package docs

// This file is a placeholder for generated swagger docs
// Run swag init in the root directory to generate proper documentation

// @title Quiz API
// @version 1.0
// @description This is a Quiz application API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.nosmo.com/support
// @contact.email support@nosmo.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description API key authentication

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT authorization with Bearer prefix

// SwaggerInfo holds the API information used by the OpenAPI spec
var SwaggerInfo = struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}{
	Version:     "1.0",
	Host:        "localhost:8080",
	BasePath:    "/api/v1",
	Schemes:     []string{"http", "https"},
	Title:       "Quiz API",
	Description: "This is a Quiz application API",
}
