package main

import (
	"app/pkg/quiz/domain/entity"
	"log"

	"gorm.io/gorm"
)

// setupQuizDatabase initializes the quiz database and runs migrations
func setupQuizDatabase(db *gorm.DB) error {
	// Auto migrate quiz entities
	err := db.AutoMigrate(
		// Quiz entities
		&entity.Quiz{},
		&entity.Question{},
		&entity.Option{},
		&entity.Answer{},
		&entity.Submission{},
	)
	if err != nil {
		log.Printf("Failed to migrate quiz database: %v", err)
		return err
	}

	log.Println("Quiz database migration completed successfully")
	return nil
}
