package main

import (
	articleHandlers "app/pkg/article/handlers"
	"app/pkg/article/repository"
	articleServices "app/pkg/article/services"
	"app/pkg/middleware"
	"app/pkg/quiz/config"
	"app/pkg/quiz/handlers"
	authMiddleware "app/pkg/quiz/middleware"
	quizRepository "app/pkg/quiz/repository/gorm"
	quizServices "app/pkg/quiz/services"
	"app/pkg/validation"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/ilyakaznacheev/cleanenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "app/cmd/quiz/docs" // Import generated swagger docs
	"app/cmd/quiz/migrations"
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

// @host 157.245.61.194:8082
// @BasePath /api
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
	var cfg config.QuizConfig
	var configPath = flag.String("config", filepath.Join("cmd", "quiz", "config", "config.yml"), "path to config file")

	flag.Parse()

	// Load configuration
	if err := cleanenv.ReadConfig(*configPath, &cfg); err != nil {
		if os.IsNotExist(err) {
			if err := cleanenv.ReadEnv(&cfg); err != nil {
				log.Fatalf("error reading environment variables: %v", err)
			}
		}
		log.Fatalf("error reading config file: %v", err)
	}

	// Initialize database connection
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup quiz database (migrations)
	if err := migrations.SetupQuizDatabase(db); err != nil {
		log.Fatalf("Failed to setup quiz database: %v", err)
	}

	// Setup article database (migrations)
	if err := migrations.SetupArticleDatabase(db); err != nil {
		log.Fatalf("Failed to setup article database: %v", err)
	}

	// Initialize quiz repositories
	quizRepo := quizRepository.NewQuizRepository(db)
	questionRepo := quizRepository.NewQuestionRepository(db)
	answerRepo := quizRepository.NewAnswerRepository(db)
	submissionRepo := quizRepository.NewSubmissionRepository(db)
	userRepo := quizRepository.NewUserRepository(db)
	fileRepo := repository.NewFileRepository(db, cfg.App.UploadPath)

	// Initialize article repositories
	articleRepo := repository.NewArticleRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)

	// Initialize quiz services
	quizService := quizServices.NewQuizService(quizRepo)
	questionService := quizServices.NewQuestionService(questionRepo)
	answerService := quizServices.NewAnswerService(answerRepo)
	submissionService := quizServices.NewSubmissionService(submissionRepo)
	userService := quizServices.NewUserService(userRepo)
	authService := quizServices.NewAuthService(userRepo, &cfg.Auth)

	// Initialize article services
	articleService := articleServices.NewArticleService(articleRepo, categoryRepo, tagRepo, fileRepo)
	categoryService := articleServices.NewCategoryService(categoryRepo)
	tagService := articleServices.NewTagService(tagRepo, articleRepo)

	// Initialize middlewares
	errorMiddleware := middleware.NewErrorMiddleware()
	authMiddleware := authMiddleware.NewAuthMiddleware(authService)

	// Register Plugins
	validation := validation.NewValidation()

	// Register quiz routes
	quizHandler := handlers.NewQuizHandler(quizService, validation)
	questionHandler := handlers.NewQuestionHandler(questionService, validation)
	answerHandler := handlers.NewAnswerHandler(answerService, validation)
	submissionHandler := handlers.NewSubmissionHandler(submissionService, validation)
	userHandler := handlers.NewUserHandler(userService, validation)
	authHandler := handlers.NewAuthHandler(authService, validation)

	// Register article routes
	articleHandler := articleHandlers.NewArticleHandler(articleService, fileRepo, validation)
	categoryHandler := articleHandlers.NewCategoryHandler(categoryService, validation)
	tagHandler := articleHandlers.NewTagHandler(tagService, validation)

	// Initialize Fiber app with error handler from config
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
	app.Get("/docs/*", swagger.HandlerDefault)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":      "ok",
			"time":        time.Now().Format(time.RFC3339),
			"service":     "quiz",
			"version":     "1.0.0",
			"environment": cfg.App.Environment,
		})
	})

	api := app.Group("/api")

	// Register article routes
	articleHandler.RegisterRoutes(api, authMiddleware.RequireAdmin())
	categoryHandler.RegisterRoutes(api, authMiddleware.RequireAdmin())
	tagHandler.RegisterRoutes(api, authMiddleware.RequireAdmin())

	// Register quiz routes
	quizHandler.RegisterRoutes(api, authMiddleware)
	questionHandler.RegisterRoutes(api, authMiddleware)
	answerHandler.RegisterRoutes(api, authMiddleware)
	submissionHandler.RegisterRoutes(api, authMiddleware)
	userHandler.RegisterRoutes(api, authMiddleware)
	authHandler.RegisterRoutes(api)

	// Start server
	go func() {
		log.Printf("Server started on port %s", cfg.Server.Port)
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			log.Printf("Server error: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
	}
}
