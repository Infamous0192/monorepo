package question

import (
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/domain/repository"
	"app/pkg/types/pagination"
	"context"
	"slices"
)

type questionService struct {
	questionRepo repository.QuestionRepository
	answerRepo   repository.AnswerRepository
}

// NewQuestionService creates a new instance of QuestionService
func NewQuestionService(
	questionRepo repository.QuestionRepository,
	answerRepo repository.AnswerRepository,
) QuestionService {
	return &questionService{
		questionRepo: questionRepo,
		answerRepo:   answerRepo,
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
		Text:   questionDTO.Text,
		QuizID: questionDTO.QuizID,
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

	question.Text = questionDTO.Text
	question.QuizID = questionDTO.QuizID

	// Process answers - updating existing ones and adding new ones
	if len(questionDTO.Answers) > 0 {
		// Create a map of existing answers by ID for quick lookup
		existingAnswers := make(map[uint]entity.Answer)
		var answersToKeep []uint

		for _, answer := range question.Answers {
			existingAnswers[answer.ID] = answer
		}

		// Process each answer from DTO
		updatedAnswers := make([]entity.Answer, 0, len(questionDTO.Answers))

		for _, answerDTO := range questionDTO.Answers {
			// If answer has ID and exists, update it
			if answerDTO.ID != nil && *answerDTO.ID > 0 {
				if _, exists := existingAnswers[*answerDTO.ID]; exists {
					// Update the answer using the answer repository
					updatedAnswer := &entity.Answer{
						ID:         *answerDTO.ID,
						Text:       answerDTO.Text,
						Value:      *answerDTO.Value,
						QuestionID: id,
					}
					if err := s.answerRepo.Update(ctx, updatedAnswer); err != nil {
						return nil, err
					}

					updatedAnswers = append(updatedAnswers, *updatedAnswer)
					answersToKeep = append(answersToKeep, *answerDTO.ID)
				} else {
					// ID provided but not found - create new with specified ID
					newAnswer := entity.Answer{
						ID:         *answerDTO.ID,
						Text:       answerDTO.Text,
						Value:      *answerDTO.Value,
						QuestionID: id,
					}

					if err := s.answerRepo.Create(ctx, &newAnswer); err != nil {
						return nil, err
					}

					updatedAnswers = append(updatedAnswers, newAnswer)
					answersToKeep = append(answersToKeep, *answerDTO.ID)
				}
			} else {
				// No ID provided, create new
				newAnswer := entity.Answer{
					Text:       answerDTO.Text,
					Value:      *answerDTO.Value,
					QuestionID: id,
				}

				if err := s.answerRepo.Create(ctx, &newAnswer); err != nil {
					return nil, err
				}

				updatedAnswers = append(updatedAnswers, newAnswer)
			}
		}

		// Delete answers that are not in the updated list
		for answerID := range existingAnswers {
			shouldDelete := !slices.Contains(answersToKeep, answerID)

			if shouldDelete {
				if err := s.answerRepo.Delete(ctx, answerID); err != nil {
					return nil, err
				}
			}
		}

		question.Answers = updatedAnswers
	} else {
		// If no answers provided, delete all existing by finding answers with this question ID
		for _, answer := range question.Answers {
			if err := s.answerRepo.Delete(ctx, answer.ID); err != nil {
				return nil, err
			}
		}
		question.Answers = []entity.Answer{} // Set empty slice to avoid nil
	}

	// Update question data without affecting answers
	if err := s.questionRepo.Update(ctx, question); err != nil {
		return nil, err
	}

	return question, nil
}

// Delete removes a question
func (s *questionService) Delete(ctx context.Context, id uint) error {
	return s.questionRepo.Delete(ctx, id)
}
