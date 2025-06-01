package utils

import (
	"log"

	"gorm.io/gorm"
)

func (app *models.App) SaveNewArticles(articles []Article) ([]Article, error) {
	var newArticles []Article

	for _, article := range articles {
		// Check if article already exists (by content hash)
		var existing Article
		result := app.DB.Where("content_hash = ?", article.ContentHash).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			// Article is new, save it
			if err := app.DB.Create(&article).Error; err != nil {
				log.Printf("Error saving article: %v", err)
				continue
			}
			newArticles = append(newArticles, article)
			log.Printf("NEW ARTICLE: %s", article.Title)
		}
		// If article exists, skip it
	}

	return newArticles, nil
}
