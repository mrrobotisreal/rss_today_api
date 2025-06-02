package db

import (
	"fmt"
	"log"
	"os"

	"github.com/mrrobotisreal/rss_today_api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(app *models.App) error {
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
	AddDefaultSources(app)

	log.Println("Database initialized successfully")
	return nil
}

func AddDefaultSources(app *models.App) {
	defaultSources := []models.NewsSource{
		// {
		// 	Name:   "Google News",
		// 	URL:    "https://news.google.com/",
		// 	RSSURL: "https://news.google.com/rss",
		// 	Active: true,
		// },
		{
			Name:   "Fox News",
			URL:    "https://www.foxnews.com/",
			RSSURL: "https://moxie.foxnews.com/google-publisher/world.xml",
			Active: true,
		},
		{
			Name:   "Sky News",
			URL:    "https://news.sky.com/",
			RSSURL: "https://feeds.skynews.com/feeds/rss/world.xml",
			Active: true,
		},
		// {
		// 	Name:   "Breitbart",
		// 	URL:    "https://www.breitbart.com/",
		// 	RSSURL: "https://feeds.feedburner.com/breitbart",
		// 	Active: true,
		// },
		{
			Name:   "New York Post",
			URL:    "https://nypost.com/",
			RSSURL: "https://nypost.com/feed/",
			Active: true,
		},
		{
			Name:   "BBC News",
			URL:    "https://www.bbc.com/news",
			RSSURL: "http://feeds.bbci.co.uk/news/rss.xml",
			Active: true,
		},
		{
			Name:   "Reuters",
			URL:    "https://www.reuters.com",
			RSSURL: "https://feeds.reuters.com/Reuters/worldNews",
			Active: true,
		},
		{
			Name:   "CNN",
			URL:    "https://www.cnn.com",
			RSSURL: "http://rss.cnn.com/rss/edition.rss",
			Active: true,
		},
	}

	for _, source := range defaultSources {
		var existingSource models.NewsSource
		result := app.DB.Where("rss_url = ?", source.RSSURL).First(&existingSource)
		if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
			if err := app.DB.Create(&source).Error; err != nil {
				log.Printf("Error creating default source %s: %v", source.Name, err)
			} else {
				log.Printf("Added default source: %s", source.Name)
			}
		}
	}
}
