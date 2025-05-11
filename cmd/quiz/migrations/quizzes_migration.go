package migrations

import (
	"app/pkg/quiz/domain/entity"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SetupQuizDatabase initializes the quiz database and runs migrations
func SetupQuizDatabase(db *gorm.DB) error {
	// Auto migrate quiz entities
	err := db.AutoMigrate(
		// Quiz entities
		&entity.Quiz{},
		&entity.Question{},
		&entity.Answer{},
		&entity.Submission{},
		&entity.User{},
	)
	if err != nil {
		log.Printf("Failed to migrate quiz database: %v", err)
		return err
	}

	// Seed data if tables are empty
	if err := seedQuizData(db); err != nil {
		log.Printf("Warning: Failed to seed quiz data: %v", err)
	}

	// Seed user data if table is empty
	if err := seedUserData(db); err != nil {
		log.Printf("Warning: Failed to seed user data: %v", err)
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
			Text:      "What does HTML stand for?",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quiz.ID,
			Text:      "Which language is used for styling web pages?",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quiz.ID,
			Text:      "Which of the following are JavaScript frameworks?",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&questions).Error; err != nil {
		return err
	}

	// Create options for each question
	options := []entity.Answer{
		// Options for question 1
		{
			QuestionID: questions[0].ID,
			Text:       "Hyper Text Markup Language",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: questions[0].ID,
			Text:       "High Tech Multi Language",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: questions[0].ID,
			Text:       "Hyper Transfer Markup Language",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		// Options for question 2
		{
			QuestionID: questions[1].ID,
			Text:       "HTML",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: questions[1].ID,
			Text:       "CSS",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: questions[1].ID,
			Text:       "JavaScript",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		// Options for question 3
		{
			QuestionID: questions[2].ID,
			Text:       "React",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: questions[2].ID,
			Text:       "Angular",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: questions[2].ID,
			Text:       "Python",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: questions[2].ID,
			Text:       "Vue",
			Value:      0,
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

// seedUserData populates the database with initial user data if the table is empty
func seedUserData(db *gorm.DB) error {
	// Check if we have any users
	var count int64
	if err := db.Model(&entity.User{}).Count(&count).Error; err != nil {
		return err
	}

	// If users already exist, skip seeding
	if count > 0 {
		log.Println("User data already exists, skipping seed")
		return nil
	}

	log.Println("Seeding user data...")

	// Generate hashed passwords
	adminPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userPassword, err := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create sample users
	users := []entity.User{
		{
			Name:      "Admin User",
			Username:  "admin",
			Password:  string(adminPassword),
			Role:      entity.RoleAdmin,
			Status:    true,
			BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Regular User",
			Username:  "user",
			Password:  string(userPassword),
			Role:      entity.RoleUser,
			Status:    true,
			BirthDate: time.Date(1995, 5, 15, 0, 0, 0, 0, time.UTC),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Inactive User",
			Username:  "inactive",
			Password:  string(userPassword),
			Role:      entity.RoleUser,
			Status:    false,
			BirthDate: time.Date(1992, 8, 20, 0, 0, 0, 0, time.UTC),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&users).Error; err != nil {
		return err
	}

	log.Println("User data seeded successfully")
	return nil
}
