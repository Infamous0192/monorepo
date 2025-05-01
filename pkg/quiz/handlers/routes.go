package handlers

import (
	"app/pkg/quiz/middleware"
	"app/pkg/quiz/services"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes registers all quiz-related routes
func SetupRoutes(app *fiber.App, services *services.Services) {
	// Create handlers
	quizHandler := NewQuizHandler(services)
	questionHandler := NewQuestionHandler(services)
	answerHandler := NewAnswerHandler(services)
	submissionHandler := NewSubmissionHandler(services)
	authHandler := NewAuthHandler(services.Auth)
	userHandler := NewUserHandler(services)

	// Create middlewares
	authMiddleware := middleware.NewAuthMiddleware(services.Auth)

	// API group with prefix
	api := app.Group("/api/v1")

	// Authentication routes (public)
	authRoutes := api.Group("/auth")
	authRoutes.Post("/register", authHandler.Register)
	authRoutes.Post("/login", authHandler.Login)
	authRoutes.Get("/profile", authMiddleware.RequireAuth(), authHandler.GetProfile)

	// Public routes (no authentication required)
	// GET quizzes and questions are public
	api.Get("/quizzes", quizHandler.GetQuizzes)
	api.Get("/quizzes/:id", quizHandler.GetQuiz)
	api.Get("/questions", questionHandler.GetQuestions)
	api.Get("/questions/:id", questionHandler.GetQuestion)

	// User management routes - Admin only
	userRoutes := api.Group("/users", authMiddleware.RequireAdmin())
	userRoutes.Get("/", userHandler.GetUsers)
	userRoutes.Get("/:id", userHandler.GetUser)
	userRoutes.Post("/", userHandler.CreateUser)
	userRoutes.Put("/:id", userHandler.UpdateUser)
	userRoutes.Delete("/:id", userHandler.DeleteUser)

	// Protected routes (authentication required)
	// For quiz management - Admin only
	quizRoutes := api.Group("/quizzes", authMiddleware.RequireAdmin())
	quizRoutes.Post("/", quizHandler.CreateQuiz)
	quizRoutes.Put("/:id", quizHandler.UpdateQuiz)
	quizRoutes.Delete("/:id", quizHandler.DeleteQuiz)

	// Question management - Admin only
	questionRoutes := api.Group("/questions", authMiddleware.RequireAdmin())
	questionRoutes.Post("/", questionHandler.CreateQuestion)
	questionRoutes.Put("/:id", questionHandler.UpdateQuestion)
	questionRoutes.Delete("/:id", questionHandler.DeleteQuestion)

	// Answer management - Admin only
	answerRoutes := api.Group("/answers", authMiddleware.RequireAdmin())
	answerRoutes.Get("/", answerHandler.GetAnswers)
	answerRoutes.Get("/:id", answerHandler.GetAnswer)
	answerRoutes.Post("/", answerHandler.CreateAnswer)
	answerRoutes.Put("/:id", answerHandler.UpdateAnswer)
	answerRoutes.Delete("/:id", answerHandler.DeleteAnswer)

	// Submission routes - Authenticated only
	submissionRoutes := api.Group("/submissions", authMiddleware.RequireAuth())
	submissionRoutes.Get("/", submissionHandler.GetSubmissions)
	submissionRoutes.Get("/:id", submissionHandler.GetSubmission)           // User can see their own submission
	submissionRoutes.Post("/", submissionHandler.CreateSubmission)          // Any authenticated user can submit
	submissionRoutes.Post("/bulk", submissionHandler.CreateBulkSubmissions) // Any authenticated user can submit in bulk

	// Admin-only submission operations
	adminSubmissionRoutes := api.Group("/submissions", authMiddleware.RequireAdmin())
	adminSubmissionRoutes.Put("/:id", submissionHandler.UpdateSubmission)
	adminSubmissionRoutes.Delete("/:id", submissionHandler.DeleteSubmission)
}
