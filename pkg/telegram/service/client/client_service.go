package client

import (
	"app/pkg/telegram/domain/entity"
	"app/pkg/telegram/domain/repository"
	"app/pkg/types/pagination"
	"context"
)

type clientService struct {
	clientRepo repository.ClientRepository
}

// NewClientService creates a new instance of ClientService
func NewClientService(clientRepo repository.ClientRepository) ClientService {
	return &clientService{
		clientRepo: clientRepo,
	}
}

// Get retrieves a bot client by ID
func (s *clientService) Get(ctx context.Context, id string) (*entity.Client, error) {
	return s.clientRepo.Get(ctx, id)
}

// GetByToken retrieves a bot client by token
func (s *clientService) GetByToken(ctx context.Context, token string) (*entity.Client, error) {
	client, err := s.clientRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetAll retrieves multiple bot clients with pagination
func (s *clientService) GetAll(ctx context.Context, pag pagination.Pagination) ([]*entity.Client, int64, error) {
	return s.clientRepo.GetAll(ctx, pag)
}

// CreateClient creates a new client
func (s *clientService) Create(ctx context.Context, client *entity.Client) error {
	return s.clientRepo.Create(ctx, client)
}

// UpdateClient modifies an existing client
func (s *clientService) Update(ctx context.Context, client *entity.Client) error {
	return s.clientRepo.Update(ctx, client)
}

// DeleteClient removes a client
func (s *clientService) Delete(ctx context.Context, id string) error {
	return s.clientRepo.Delete(ctx, id)
}
