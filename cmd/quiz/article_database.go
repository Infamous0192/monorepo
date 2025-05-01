package main

import (
	"app/pkg/article/domain/entity"
	"log"

	"gorm.io/gorm"
)

// setupArticleDatabase initializes the article database and runs migrations
func setupArticleDatabase(db *gorm.DB) error {
	// Auto migrate article entities
	err := db.AutoMigrate(
		// Article entities
		&entity.Article{},
		&entity.Category{},
		&entity.Tag{},
	)
	if err != nil {
		log.Printf("Failed to migrate article database: %v", err)
		return err
	}

	log.Println("Article database migration completed successfully")
	return nil
}
