package main

import (
	articleHandlers "app/pkg/article/handlers"
	"app/pkg/article/repository"
	articleServices "app/pkg/article/services"
	"app/pkg/middleware"
	"app/pkg/quiz/handlers"
	quizRepository "app/pkg/quiz/repository/gorm"
	"app/pkg/quiz/services"
	"app/pkg/validation"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "app/cmd/quiz/docs" // Import generated swagger docs
)

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

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database connection
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "quiz_db"),
		getEnv("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup quiz database (migrations)
	if err := setupQuizDatabase(db); err != nil {
		log.Fatalf("Failed to setup quiz database: %v", err)
	}

	// Setup article database (migrations)
	if err := setupArticleDatabase(db); err != nil {
		log.Fatalf("Failed to setup article database: %v", err)
	}

	// Initialize error middleware
	errorMiddleware := middleware.NewErrorMiddleware()
	keyMiddleware := middleware.NewKeyMiddleware(getEnv("API_KEY", "default-api-key"))

	// Initialize validation
	validation.NewValidation() // Initialize validation package

	// Initialize quiz repositories
	quizRepo := quizRepository.NewQuizRepository(db)
	questionRepo := quizRepository.NewQuestionRepository(db)
	answerRepo := quizRepository.NewAnswerRepository(db)
	submissionRepo := quizRepository.NewSubmissionRepository(db)
	userRepo := quizRepository.NewUserRepository(db)

	// Initialize article repositories
	articleRepo := repository.NewArticleRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)

	// Initialize quiz services
	quizServices := services.NewServices(
		db,
		quizRepo,
		questionRepo,
		answerRepo,
		submissionRepo,
		userRepo,
	)

	// Initialize article services
	articleSvc := articleServices.NewArticleService(articleRepo, categoryRepo, tagRepo)
	categorySvc := articleServices.NewCategoryService(categoryRepo)
	tagSvc := articleServices.NewTagService(tagRepo, articleRepo)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: errorMiddleware.Handler(),
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-API-Key",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Register quiz routes
	handlers.SetupRoutes(app, quizServices)

	// Register article routes
	articleHandler := articleHandlers.NewArticleHandler(articleSvc, categorySvc, tagSvc)
	categoryHandler := articleHandlers.NewCategoryHandler(categorySvc)
	tagHandler := articleHandlers.NewTagHandler(tagSvc)

	// Set up routes for each handler
	articleHandler.RegisterRoutes(app, keyMiddleware.ValidateKey())
	categoryHandler.RegisterRoutes(app, keyMiddleware.ValidateKey())
	tagHandler.RegisterRoutes(app, keyMiddleware.ValidateKey())

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server started on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnv gets an environment variable or returns the default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
