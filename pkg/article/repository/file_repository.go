package repository

import (
	"app/pkg/article/domain/entity"
	"app/pkg/article/domain/repository"
	"app/pkg/exception"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

// fileRepository implements the FileRepository interface
type fileRepository struct {
	db         *gorm.DB
	uploadPath string
}

// NewFileRepository creates a new instance of the file repository
func NewFileRepository(db *gorm.DB, uploadPath string) repository.FileRepository {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}

	return &fileRepository{
		db:         db,
		uploadPath: uploadPath,
	}
}

// Store saves a file and returns the file entity
func (r *fileRepository) Store(ctx context.Context, name string, contentType string, size int64, reader io.Reader) (*entity.File, error) {
	// Generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), name)
	filePath := filepath.Join(r.uploadPath, filename)

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return nil, exception.InternalError(fmt.Sprintf("Failed to create file: %v", err))
	}
	defer file.Close()

	// Copy file content
	if _, err := io.Copy(file, reader); err != nil {
		os.Remove(filePath) // Clean up on error
		return nil, exception.InternalError(fmt.Sprintf("Failed to write file: %v", err))
	}

	// Create file entity
	fileEntity := &entity.File{
		Name:        name,
		Path:        filename,
		Size:        size,
		ContentType: contentType,
	}

	if err := r.db.WithContext(ctx).Create(fileEntity).Error; err != nil {
		os.Remove(filePath) // Clean up on error
		return nil, exception.InternalError(fmt.Sprintf("Failed to store file metadata: %v", err))
	}

	return fileEntity, nil
}

// FindOne retrieves a single file by ID
func (r *fileRepository) FindOne(ctx context.Context, id uint) (*entity.File, error) {
	var file entity.File

	result := r.db.WithContext(ctx).First(&file, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, exception.NotFound("File")
		}
		return nil, exception.InternalError(fmt.Sprintf("Failed to find file: %v", result.Error))
	}

	return &file, nil
}

// Delete removes a file
func (r *fileRepository) Delete(ctx context.Context, id uint) error {
	// Get file info first
	file, err := r.FindOne(ctx, id)
	if err != nil {
		return err
	}

	// Delete physical file
	filePath := filepath.Join(r.uploadPath, file.Path)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return exception.InternalError(fmt.Sprintf("Failed to delete physical file: %v", err))
	}

	// Delete from database
	result := r.db.WithContext(ctx).Delete(&entity.File{}, id)
	if result.Error != nil {
		return exception.InternalError(fmt.Sprintf("Failed to delete file record: %v", result.Error))
	}

	return nil
}
