package quiz

import (
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/domain/repository"
	"app/pkg/types/pagination"
	"context"
	"fmt"
)

type quizService struct {
	quizRepo repository.QuizRepository
}

// NewQuizService creates a new instance of QuizService
func NewQuizService(quizRepo repository.QuizRepository) QuizService {
	return &quizService{
		quizRepo: quizRepo,
	}
}

// FindOne retrieves a quiz by ID
func (s *quizService) FindOne(ctx context.Context, id uint) (*entity.Quiz, error) {
	return s.quizRepo.FindOne(ctx, id)
}

// FindAll retrieves multiple quizzes with filtering and pagination
func (s *quizService) FindAll(ctx context.Context, query entity.QuizQuery) (*pagination.PaginatedResult[entity.Quiz], error) {
	items, total, err := s.quizRepo.FindAll(ctx, query)
	if err != nil {
		return nil, err
	}

	// Calculate pagination metadata
	hasNext := (int64(query.Page*query.GetLimit()) < total)
	hasPrev := query.Page > 1

	result := &pagination.PaginatedResult[entity.Quiz]{
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

// Create creates a new quiz
func (s *quizService) Create(ctx context.Context, quizDTO entity.QuizDTO) (*entity.Quiz, error) {
	// Log incoming DTO
	fmt.Printf("Creating quiz with DTO: %+v\n", quizDTO)

	quiz := &entity.Quiz{
		Name:        quizDTO.Name,
		Description: quizDTO.Description,
	}

	// Log quiz entity before creation
	fmt.Printf("Quiz entity before creation: %+v\n", quiz)

	if err := s.quizRepo.Create(ctx, quiz); err != nil {
		fmt.Printf("Error creating quiz: %v\n", err)
		return nil, err
	}

	fmt.Printf("Quiz created successfully: %+v\n", quiz)
	return quiz, nil
}

// Update modifies an existing quiz
func (s *quizService) Update(ctx context.Context, id uint, quizDTO entity.QuizDTO) (*entity.Quiz, error) {
	quiz, err := s.quizRepo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	if quiz == nil {
		return nil, nil
	}

	quiz.Name = quizDTO.Name
	quiz.Description = quizDTO.Description

	if err := s.quizRepo.Update(ctx, quiz); err != nil {
		return nil, err
	}

	return quiz, nil
}

// Delete removes a quiz
func (s *quizService) Delete(ctx context.Context, id uint) error {
	return s.quizRepo.Delete(ctx, id)
}
