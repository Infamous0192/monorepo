package repository

import (
	"app/pkg/article/domain/entity"
	"context"
	"io"
)

// FileRepository defines the interface for file data access
type FileRepository interface {
	// Store saves a file and returns the file entity
	Store(ctx context.Context, name string, contentType string, size int64, reader io.Reader) (*entity.File, error)

	// FindOne retrieves a single file by ID
	FindOne(ctx context.Context, id uint) (*entity.File, error)

	// Delete removes a file
	Delete(ctx context.Context, id uint) error
}
