package db

import (
	"fmt"
	"log"
	"os"

	"github.com/mrrobotisreal/rss_today_api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (app *models.App) InitDatabase() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	app.DB = db

	// Create all tables
	err = db.AutoMigrate(&models.User{}, &models.NewsSource{}, &models.Article{}, &models.UserAlert{}, &models.NotificationSent{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	// Add default news sources if they don't exist
	app.addDefaultSources()

	log.Println("Database initialized successfully")
	return nil
}
