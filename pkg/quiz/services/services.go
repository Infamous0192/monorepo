package services

import (
	"app/pkg/quiz/config"
	"app/pkg/quiz/domain/repository"
	"app/pkg/quiz/services/answer"
	"app/pkg/quiz/services/auth"
	"app/pkg/quiz/services/question"
	"app/pkg/quiz/services/quiz"
	"app/pkg/quiz/services/submission"
	"app/pkg/quiz/services/user"

	"gorm.io/gorm"
)

// Services contains all application services
type Services struct {
	Quiz       quiz.QuizService
	Question   question.QuestionService
	Answer     answer.AnswerService
	Submission submission.SubmissionService
	Auth       auth.AuthService
	User       user.UserService
}

// NewServices creates and returns all application services
func NewServices(
	db *gorm.DB,
	quizRepo repository.QuizRepository,
	questionRepo repository.QuestionRepository,
	answerRepo repository.AnswerRepository,
	submissionRepo repository.SubmissionRepository,
	userRepo repository.UserRepository,
) *Services {
	// Create auth config
	authConfig := config.NewAuthConfig()

	return &Services{
		Quiz:       quiz.NewQuizService(quizRepo),
		Question:   question.NewQuestionService(questionRepo),
		Answer:     answer.NewAnswerService(answerRepo),
		Submission: submission.NewSubmissionService(submissionRepo),
		Auth:       auth.NewAuthService(userRepo, authConfig),
		User:       user.NewUserService(userRepo),
	}
}
