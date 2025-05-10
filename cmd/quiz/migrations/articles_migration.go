package migrations

import (
	"app/pkg/article/domain/entity"
	"log"
	"time"

	"gorm.io/gorm"
)

// SetupArticleDatabase initializes the article database and runs migrations
func SetupArticleDatabase(db *gorm.DB) error {
	// Auto migrate article entities
	err := db.AutoMigrate(
		// Article entities
		&entity.Article{},
		&entity.Category{},
		&entity.Tag{},
		&entity.File{},
	)
	if err != nil {
		log.Printf("Failed to migrate article database: %v", err)
		return err
	}

	// Seed data if tables are empty
	if err := seedArticleData(db); err != nil {
		log.Printf("Warning: Failed to seed article data: %v", err)
	}

	log.Println("Article database migration completed successfully")
	return nil
}

// seedArticleData populates the database with initial article data if tables are empty
func seedArticleData(db *gorm.DB) error {
	// Check if we have any categories, tags, or articles
	var categoryCount, tagCount, articleCount int64

	if err := db.Model(&entity.Category{}).Count(&categoryCount).Error; err != nil {
		return err
	}

	if err := db.Model(&entity.Tag{}).Count(&tagCount).Error; err != nil {
		return err
	}

	if err := db.Model(&entity.Article{}).Count(&articleCount).Error; err != nil {
		return err
	}

	// If all entities already exist, skip seeding
	if categoryCount > 0 && tagCount > 0 && articleCount > 0 {
		log.Println("Article data already exists, skipping seed")
		return nil
	}

	log.Println("Seeding article data...")

	// Create categories
	categories := []entity.Category{
		{
			Name:        "Programming",
			Description: "Articles about programming topics",
			Slug:        "programming",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Web Development",
			Description: "Articles about web development",
			Slug:        "web-development",
			ParentID:    nil, // Will be updated after creation
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Mobile Development",
			Description: "Articles about mobile app development",
			Slug:        "mobile-development",
			ParentID:    nil, // Will be updated after creation
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if err := db.Create(&categories).Error; err != nil {
		return err
	}

	// Update parent relationships for categories
	categories[1].ParentID = &categories[0].ID // Web Development -> Programming
	categories[2].ParentID = &categories[0].ID // Mobile Development -> Programming

	if err := db.Save(&categories[1]).Error; err != nil {
		return err
	}
	if err := db.Save(&categories[2]).Error; err != nil {
		return err
	}

	// Create tags
	tags := []entity.Tag{
		{
			Name:        "Go",
			Description: "Articles related to Go programming language",
			Slug:        "go",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "JavaScript",
			Description: "Articles related to JavaScript",
			Slug:        "javascript",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "React",
			Description: "Articles related to React.js",
			Slug:        "react",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "API",
			Description: "Articles related to API development",
			Slug:        "api",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if err := db.Create(&tags).Error; err != nil {
		return err
	}

	// Create articles
	now := time.Now()
	publishedAt := now.Add(-24 * time.Hour) // Yesterday

	articles := []entity.Article{
		{
			Title:       "Getting Started with Go Fiber",
			Content:     "<h1>Getting Started with Go Fiber</h1><p>Go Fiber is a web framework built on top of Fasthttp, the fastest HTTP engine for Go. This tutorial will guide you through setting up your first Fiber API.</p>",
			Slug:        "getting-started-with-go-fiber",
			PublishedAt: &publishedAt,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "Building Modern Web Apps with React",
			Content:     "<h1>Building Modern Web Apps with React</h1><p>React is a popular JavaScript library for building user interfaces. In this article, we'll explore how to create modern web applications using React.</p>",
			Slug:        "building-modern-web-apps-with-react",
			PublishedAt: &publishedAt,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "RESTful API Design Principles",
			Content:     "<h1>RESTful API Design Principles</h1><p>Designing a good RESTful API is crucial for creating scalable and maintainable applications. This article covers the key principles to follow.</p>",
			Slug:        "restful-api-design-principles",
			PublishedAt: nil, // Draft article
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	if err := db.Create(&articles).Error; err != nil {
		return err
	}

	// Set up article-category relationships
	if err := db.Exec("INSERT INTO article_categories (article_id, category_id) VALUES (?, ?)", articles[0].ID, categories[1].ID).Error; err != nil {
		return err
	}
	if err := db.Exec("INSERT INTO article_categories (article_id, category_id) VALUES (?, ?)", articles[1].ID, categories[1].ID).Error; err != nil {
		return err
	}
	if err := db.Exec("INSERT INTO article_categories (article_id, category_id) VALUES (?, ?)", articles[2].ID, categories[0].ID).Error; err != nil {
		return err
	}

	// Set up article-tag relationships
	if err := db.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articles[0].ID, tags[0].ID).Error; err != nil {
		return err
	}
	if err := db.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articles[0].ID, tags[3].ID).Error; err != nil {
		return err
	}
	if err := db.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articles[1].ID, tags[1].ID).Error; err != nil {
		return err
	}
	if err := db.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articles[1].ID, tags[2].ID).Error; err != nil {
		return err
	}
	if err := db.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articles[2].ID, tags[0].ID).Error; err != nil {
		return err
	}
	if err := db.Exec("INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)", articles[2].ID, tags[3].ID).Error; err != nil {
		return err
	}

	log.Println("Article data seeded successfully")
	return nil
}
