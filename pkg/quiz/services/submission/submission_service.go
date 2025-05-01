package submission

import (
	"app/pkg/quiz/domain/entity"
	"app/pkg/quiz/domain/repository"
	"app/pkg/types/pagination"
	"context"
)

type submissionService struct {
	submissionRepo repository.SubmissionRepository
}

// NewSubmissionService creates a new instance of SubmissionService
func NewSubmissionService(submissionRepo repository.SubmissionRepository) SubmissionService {
	return &submissionService{
		submissionRepo: submissionRepo,
	}
}

// FindOne retrieves a submission by ID
func (s *submissionService) FindOne(ctx context.Context, id uint) (*entity.Submission, error) {
	return s.submissionRepo.FindOne(ctx, id)
}

// FindAll retrieves multiple submissions with filtering and pagination
func (s *submissionService) FindAll(ctx context.Context, query entity.SubmissionQuery) (*pagination.PaginatedResult[entity.Submission], error) {
	items, total, err := s.submissionRepo.FindAll(ctx, query)
	if err != nil {
		return nil, err
	}

	// Calculate pagination metadata
	hasNext := (int64(query.Page*query.GetLimit()) < total)
	hasPrev := query.Page > 1

	result := &pagination.PaginatedResult[entity.Submission]{
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

// Create creates a new submission
func (s *submissionService) Create(ctx context.Context, submissionDTO entity.SubmissionDTO) (*entity.Submission, error) {
	submission := &entity.Submission{
		QuizID:     submissionDTO.Quiz,
		QuestionID: submissionDTO.Question,
		AnswerID:   submissionDTO.Answer,
		UserID:     submissionDTO.User,
	}

	if err := s.submissionRepo.Create(ctx, submission); err != nil {
		return nil, err
	}

	return submission, nil
}

// CreateBulk creates multiple submissions
func (s *submissionService) CreateBulk(ctx context.Context, submissionDTO entity.SubmissionInsertDTO) ([]*entity.Submission, error) {
	submissions := make([]*entity.Submission, len(submissionDTO.Data))

	for i, dto := range submissionDTO.Data {
		submissions[i] = &entity.Submission{
			QuizID:     dto.Quiz,
			QuestionID: dto.Question,
			AnswerID:   dto.Answer,
			UserID:     dto.User,
		}
	}

	if err := s.submissionRepo.CreateBulk(ctx, submissions); err != nil {
		return nil, err
	}

	return submissions, nil
}

// Update modifies an existing submission
func (s *submissionService) Update(ctx context.Context, id uint, submissionDTO entity.SubmissionDTO) (*entity.Submission, error) {
	submission, err := s.submissionRepo.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	if submission == nil {
		return nil, nil
	}

	submission.QuizID = submissionDTO.Quiz
	submission.QuestionID = submissionDTO.Question
	submission.AnswerID = submissionDTO.Answer
	submission.UserID = submissionDTO.User

	if err := s.submissionRepo.Update(ctx, submission); err != nil {
		return nil, err
	}

	return submission, nil
}

// Delete removes a submission
func (s *submissionService) Delete(ctx context.Context, id uint) error {
	return s.submissionRepo.Delete(ctx, id)
}
