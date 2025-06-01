package services

import (
	"log"

	"github.com/mrrobotisreal/rss_today_api/internal/models"
	"gorm.io/gorm"
)

func SaveNewArticles(app *models.App, articles []models.Article) ([]models.Article, error) {
	var newArticles []models.Article

	for _, article := range articles {
		// Check if article already exists by link or content hash
		var existingArticle models.Article
		result := app.DB.Where("link = ? OR content_hash = ?", article.Link, article.ContentHash).First(&existingArticle)

		if result.Error == gorm.ErrRecordNotFound {
			// Article is new, save it
			if err := app.DB.Create(&article).Error; err != nil {
				log.Printf("Error saving article '%s': %v", article.Title, err)
				continue
			}
			newArticles = append(newArticles, article)
			log.Printf("Saved new article: %s", article.Title)
		}
	}

	if len(newArticles) > 0 {
		log.Printf("Saved %d new articles to database", len(newArticles))
	}

	return newArticles, nil
}