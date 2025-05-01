package answer

import (
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/domain/repository"
	"app/pkg/types/pagination"
	"context"
)

type answerService struct {
	answerRepo repository.AnswerRepository
}

// NewAnswerService creates a new instance of AnswerService
func NewAnswerService(answerRepo repository.AnswerRepository) AnswerService {
	return &answerService{
		answerRepo: answerRepo,
	}
}

// FindOne retrieves an answer by ID
func (s *answerService) FindOne(ctx context.Context, id uint) (*entity.Answer, error) {
	return s.answerRepo.FindOne(ctx, id)
}

// FindAll retrieves multiple answers with filtering and pagination
func (s *answerService) FindAll(ctx context.Context, query entity.AnswerQuery) (*pagination.PaginatedResult[entity.Answer], error) {
	items, total, err := s.answerRepo.FindAll(ctx, query)
	if err != nil {
		return nil, err
	}

	// Calculate pagination metadata
	hasNext := (int64(query.Page*query.GetLimit()) < total)
	hasPrev := query.Page > 1

	result := &pagination.PaginatedResult[entity.Answer]{
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

// Create creates a new answer
func (s *answerService) Create(ctx context.Context, answerDTO entity.AnswerDTO) (*entity.Answer, error) {
	answer := &entity.Answer{
		Text:       answerDTO.Text,
		Value:      *answerDTO.Value,
		QuestionID: answerDTO.Question,
	}

	if err := s.answerRepo.Create(ctx, answer); err != nil {
		return nil, err
	}

	return answer, nil
}

// Update modifies an existing answer
func (s *answerService) Update(ctx context.Context, id uint, answerDTO entity.AnswerDTO) (*entity.Answer, error) {
	answer, err := s.answerRepo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	if answer == nil {
		return nil, nil
	}

	answer.Text = answerDTO.Text
	answer.Value = *answerDTO.Value

	if answerDTO.Question != 0 {
		answer.QuestionID = answerDTO.Question
	}

	if err := s.answerRepo.Update(ctx, answer); err != nil {
		return nil, err
	}

	return answer, nil
}

// Delete removes an answer
func (s *answerService) Delete(ctx context.Context, id uint) error {
	return s.answerRepo.Delete(ctx, id)
}
