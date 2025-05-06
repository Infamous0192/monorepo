package main

import (
	"app/pkg/quiz/domain/entity"
	"log"
	"time"

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

	// Seed data if tables are empty
	if err := seedQuizData(db); err != nil {
		log.Printf("Warning: Failed to seed quiz data: %v", err)
	}

	log.Println("Quiz database migration completed successfully")
	return nil
}

// seedQuizData populates the database with initial quiz data if tables are empty
func seedQuizData(db *gorm.DB) error {
	// Check if we have any quizzes
	var count int64
	if err := db.Model(&entity.Quiz{}).Count(&count).Error; err != nil {
		return err
	}

	// If quizzes already exist, skip seeding
	if count > 0 {
		log.Println("Quiz data already exists, skipping seed")
		return nil
	}

	log.Println("Seeding quiz data...")

	// Create a sample quiz
	quiz := &entity.Quiz{
		Name:        "Sample Programming Quiz",
		Description: "Test your programming knowledge with this quiz",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(quiz).Error; err != nil {
		return err
	}

	// Create questions for the quiz
	questions := []entity.Question{
		{
			QuizID:    quiz.ID,
			Content:   "What does HTML stand for?",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quiz.ID,
			Content:   "Which language is used for styling web pages?",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quiz.ID,
			Content:   "Which of the following are JavaScript frameworks?",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&questions).Error; err != nil {
		return err
	}

	// Create options for each question
	options := []entity.Option{
		// Options for question 1
		{
			QuestionID: questions[0].ID,
			Content:    "Hyper Text Markup Language",
			IsCorrect:  true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: questions[0].ID,
			Content:    "High Tech Multi Language",
			IsCorrect:  false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: questions[0].ID,
			Content:    "Hyper Transfer Markup Language",
			IsCorrect:  false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		// Options for question 2
		{
			QuestionID: questions[1].ID,
			Content:    "HTML",
			IsCorrect:  false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: questions[1].ID,
			Content:    "CSS",
			IsCorrect:  true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: questions[1].ID,
			Content:    "JavaScript",
			IsCorrect:  false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		// Options for question 3
		{
			QuestionID: questions[2].ID,
			Content:    "React",
			IsCorrect:  true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: questions[2].ID,
			Content:    "Angular",
			IsCorrect:  true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: questions[2].ID,
			Content:    "Python",
			IsCorrect:  false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: questions[2].ID,
			Content:    "Vue",
			IsCorrect:  true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&options).Error; err != nil {
		return err
	}

	log.Println("Quiz data seeded successfully")
	return nil
}
