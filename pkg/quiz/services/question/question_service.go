package question

import (
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/domain/repository"
	"app/pkg/types/pagination"
	"context"
)

type questionService struct {
	questionRepo repository.QuestionRepository
}

// NewQuestionService creates a new instance of QuestionService
func NewQuestionService(questionRepo repository.QuestionRepository) QuestionService {
	return &questionService{
		questionRepo: questionRepo,
	}
}

// FindOne retrieves a question by ID
func (s *questionService) FindOne(ctx context.Context, id uint) (*entity.Question, error) {
	return s.questionRepo.FindOne(ctx, id)
}

// FindAll retrieves multiple questions with filtering and pagination
func (s *questionService) FindAll(ctx context.Context, query entity.QuestionQuery) (*pagination.PaginatedResult[entity.Question], error) {
	items, total, err := s.questionRepo.FindAll(ctx, query)
	if err != nil {
		return nil, err
	}

	// Calculate pagination metadata
	hasNext := (int64(query.Page*query.GetLimit()) < total)
	hasPrev := query.Page > 1

	result := &pagination.PaginatedResult[entity.Question]{
		Metadata: pagination.Metadata{
			Pagination: pagination.Pagination{
				Page:  query.Page,
				Limit: query.GetLimit(),
			},
			Total:   total,
			Count:   len(items),
			HasNext: hasNext,
			HasPrev: hasPrev,
		},
		Result: items,
	}

	return result, nil
}

// Create creates a new question
func (s *questionService) Create(ctx context.Context, questionDTO entity.QuestionDTO) (*entity.Question, error) {
	question := &entity.Question{
		Content: questionDTO.Content,
		QuizID:  questionDTO.QuizID,
	}

	// Create options if provided
	if len(questionDTO.Options) > 0 {
		options := make([]entity.Option, len(questionDTO.Options))
		for i, optionDTO := range questionDTO.Options {
			options[i] = entity.Option{
				Content:   optionDTO.Content,
				IsCorrect: optionDTO.IsCorrect,
			}
		}
		question.Options = options
	}

	// Create answers if provided
	if len(questionDTO.Answers) > 0 {
		answers := make([]entity.Answer, len(questionDTO.Answers))
		for i, answerDTO := range questionDTO.Answers {
			answers[i] = entity.Answer{
				Text:  answerDTO.Text,
				Value: *answerDTO.Value,
			}
		}
		question.Answers = answers
	}

	if err := s.questionRepo.Create(ctx, question); err != nil {
		return nil, err
	}

	return question, nil
}

// Update modifies an existing question
func (s *questionService) Update(ctx context.Context, id uint, questionDTO entity.QuestionDTO) (*entity.Question, error) {
	question, err := s.questionRepo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	if question == nil {
		return nil, nil
	}

	question.Content = questionDTO.Content
	question.QuizID = questionDTO.QuizID

	// Handle options update (simplified - in real implementation would need to handle deletions)
	if len(questionDTO.Options) > 0 {
		options := make([]entity.Option, len(questionDTO.Options))
		for i, optionDTO := range questionDTO.Options {
			options[i] = entity.Option{
				Content:    optionDTO.Content,
				IsCorrect:  optionDTO.IsCorrect,
				QuestionID: id,
			}
		}
		question.Options = options
	}

	// Handle answers update (simplified)
	if len(questionDTO.Answers) > 0 {
		answers := make([]entity.Answer, len(questionDTO.Answers))
		for i, answerDTO := range questionDTO.Answers {
			answers[i] = entity.Answer{
				Text:       answerDTO.Text,
				Value:      *answerDTO.Value,
				QuestionID: id,
			}
		}
		question.Answers = answers
	}

	if err := s.questionRepo.Update(ctx, question); err != nil {
		return nil, err
	}

	return question, nil
}

// Delete removes a question
func (s *questionService) Delete(ctx context.Context, id uint) error {
	return s.questionRepo.Delete(ctx, id)
}
