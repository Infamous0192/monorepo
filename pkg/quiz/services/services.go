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
)

func NewAnswerService(answerRepo repository.AnswerRepository) answer.AnswerService {
	return answer.NewAnswerService(answerRepo)
}

func NewQuestionService(questionRepo repository.QuestionRepository, answerRepo repository.AnswerRepository) question.QuestionService {
	return question.NewQuestionService(questionRepo, answerRepo)
}

func NewQuizService(quizRepo repository.QuizRepository) quiz.QuizService {
	return quiz.NewQuizService(quizRepo)
}

func NewSubmissionService(submissionRepo repository.SubmissionRepository) submission.SubmissionService {
	return submission.NewSubmissionService(submissionRepo)
}

func NewUserService(userRepo repository.UserRepository) user.UserService {
	return user.NewUserService(userRepo)
}

func NewAuthService(userRepo repository.UserRepository, authConfig *config.AuthConfig) auth.AuthService {
	return auth.NewAuthService(userRepo, authConfig)
}
