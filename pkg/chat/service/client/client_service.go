package client

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/chat/domain/repository"
	"app/pkg/database/redis"
	"app/pkg/exception"
	"app/pkg/types/pagination"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const (
	// Cache duration for authenticated users
	authCacheDuration = 15 * time.Minute
	// Cache key prefix for auth tokens
	authCacheKeyPrefix = "auth:token:"
	// Cache duration for validated clients
	clientCacheDuration = 1 * time.Hour
	// Cache key prefix for client keys
	clientCacheKeyPrefix = "client:key:"
)

type clientService struct {
	clientRepo  repository.ClientRepository
	userRepo    repository.UserRepository
	redisClient *redis.Client
	httpClient  *http.Client
}

// NewClientService creates a new instance of ClientService
func NewClientService(clientRepo repository.ClientRepository, userRepo repository.UserRepository, redisClient *redis.Client) ClientService {
	return &clientService{
		clientRepo:  clientRepo,
		userRepo:    userRepo,
		redisClient: redisClient,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetClient retrieves a client by ID
func (s *clientService) GetClient(ctx context.Context, id string) (*entity.Client, error) {
	return s.clientRepo.Get(ctx, id)
}

// GetByKey retrieves a client by client key
func (s *clientService) GetByKey(ctx context.Context, clientKey string) (*entity.Client, error) {
	return s.clientRepo.GetByKey(ctx, clientKey)
}

// GetClients retrieves multiple clients with pagination
func (s *clientService) GetClients(ctx context.Context, pag pagination.Pagination) ([]*entity.Client, int64, error) {
	return s.clientRepo.GetAll(ctx, pag)
}

// CreateClient creates a new client
func (s *clientService) CreateClient(ctx context.Context, client *entity.Client) error {
	return s.clientRepo.Create(ctx, client)
}

// UpdateClient modifies an existing client
func (s *clientService) UpdateClient(ctx context.Context, client *entity.Client) error {
	return s.clientRepo.Update(ctx, client)
}

// DeleteClient removes a client
func (s *clientService) DeleteClient(ctx context.Context, id string) error {
	return s.clientRepo.Delete(ctx, id)
}

// Authenticate validates a JWT token against the client's auth endpoint and returns the authenticated user
func (s *clientService) Authenticate(ctx context.Context, clientID string, token string) (*entity.User, error) {
	// Try to get from cache first
	cacheKey := authCacheKeyPrefix + token
	if cachedUser, err := s.getUserFromCache(ctx, cacheKey); err == nil && cachedUser != nil {
		return cachedUser, nil
	}

	// Get client details
	client, err := s.clientRepo.Get(ctx, clientID)
	if err != nil {
		return nil, exception.InternalError("Error fetching client details")
	}
	if client == nil {
		return nil, exception.NotFound("Client")
	}

	// Create request to auth endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", client.AuthEndpoint, nil)
	if err != nil {
		return nil, exception.InternalError("Failed to create authentication request")
	}

	// Add authorization header
	req.Header.Set("Authorization", "Bearer "+token)

	// Make the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, exception.InternalError("Authentication request failed")
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return nil, exception.Http(401, "Invalid token")
		case http.StatusForbidden:
			return nil, exception.Forbidden()
		case http.StatusNotFound:
			return nil, exception.NotFound("User")
		default:
			return nil, exception.InternalError("Authentication service error: " + string(body))
		}
	}

	// Parse response body
	var user entity.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, exception.InternalError("Failed to decode authentication response")
	}

	// Check if user exists
	existingUser, err := s.userRepo.Get(ctx, user.ID)
	if err != nil {
		return nil, exception.InternalError("Error checking user existence")
	}

	// Create or update user
	if existingUser == nil {
		// Create new user
		if err := s.userRepo.Create(ctx, &user); err != nil {
			return nil, exception.InternalError("Failed to create user")
		}
	} else {
		// Update existing user
		if err := s.userRepo.Update(ctx, &user); err != nil {
			return nil, exception.InternalError("Failed to update user")
		}
	}

	// Cache the authenticated user
	if err := s.cacheUser(ctx, cacheKey, &user); err != nil {
		// Log the error but don't fail the request
		// fmt.Printf("Failed to cache user: %v\n", err)
	}

	return &user, nil
}

// getUserFromCache attempts to retrieve a user from Redis cache
func (s *clientService) getUserFromCache(ctx context.Context, key string) (*entity.User, error) {
	val, err := s.redisClient.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var user entity.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// cacheUser stores a user in Redis cache with expiration
func (s *clientService) cacheUser(ctx context.Context, key string, user *entity.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, key, string(data), authCacheDuration)
}

// ValidateKey validates a client key and returns the client if valid
func (s *clientService) ValidateKey(ctx context.Context, clientKey string) (*entity.Client, error) {
	if clientKey == "" {
		return nil, exception.Http(401, "Client key is required")
	}

	// Try to get client from cache first
	cacheKey := clientCacheKeyPrefix + clientKey
	cachedClient, err := s.getClientFromCache(ctx, cacheKey)
	if err == nil && cachedClient != nil {
		return cachedClient, nil
	}

	// If not in cache, fetch from database
	client, err := s.clientRepo.GetByKey(ctx, clientKey)
	if err != nil {
		return nil, exception.InternalError("Error fetching client")
	}
	if client == nil {
		return nil, exception.Http(401, "Invalid client key")
	}

	// Store in cache for future requests
	if err := s.storeClientInCache(ctx, cacheKey, client); err != nil {
		// Log the error but don't fail the request
		// fmt.Printf("Error storing client in cache: %v\n", err)
	}

	return client, nil
}

// getClientFromCache attempts to retrieve a client from Redis cache
func (s *clientService) getClientFromCache(ctx context.Context, key string) (*entity.Client, error) {
	val, err := s.redisClient.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var client entity.Client
	if err := json.Unmarshal([]byte(val), &client); err != nil {
		return nil, err
	}

	return &client, nil
}

// storeClientInCache stores a client in Redis cache with expiration
func (s *clientService) storeClientInCache(ctx context.Context, key string, client *entity.Client) error {
	data, err := json.Marshal(client)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, key, string(data), clientCacheDuration)
}
