package services

import (
	"log"
	"strings"
	"time"

	"github.com/mrrobotisreal/rss_today_api/internal/models"
)

func CheckAlertsForNewArticles(app *models.App, articles []models.Article) error {
	// Get all active user alerts
	var alerts []models.UserAlert
	if err := app.DB.Where("active = ?", true).Find(&alerts).Error; err != nil {
		return err
	}

	log.Printf("Checking %d alerts against %d new articles", len(alerts), len(articles))

	for _, alert := range alerts {
		matchingArticles := findMatchingArticles(articles, alert)

		if len(matchingArticles) > 0 {
			log.Printf("Alert ID %d matched %d articles", alert.ID, len(matchingArticles))

			// Here you would typically send notifications
			// For now, we'll just log and create notification records
			for _, article := range matchingArticles {
				notification := models.NotificationSent{
					UserID:    alert.UserID,
					ArticleID: article.ID,
					AlertID:   alert.ID,
					Method:    "email", // or whatever method you prefer
					SentAt:    time.Now(),
				}

				if err := app.DB.Create(&notification).Error; err != nil {
					log.Printf("Error creating notification record: %v", err)
				} else {
					log.Printf("Created notification for user %d, article: %s", alert.UserID, article.Title)
				}
			}
		}
	}

	return nil
}

func findMatchingArticles(articles []models.Article, alert models.UserAlert) []models.Article {
	var matchingArticles []models.Article

	for _, article := range articles {
		if articleMatchesAlert(article, alert) {
			matchingArticles = append(matchingArticles, article)
		}
	}

	return matchingArticles
}

func articleMatchesAlert(article models.Article, alert models.UserAlert) bool {
	// Check if any alert keywords match article keywords
	if len(alert.Keywords) > 0 {
		keywordMatch := false
		for _, alertKeyword := range alert.Keywords {
			alertKeywordLower := strings.ToLower(alertKeyword)

			// Check in article keywords
			for _, articleKeyword := range article.Keywords {
				if strings.Contains(strings.ToLower(articleKeyword), alertKeywordLower) {
					keywordMatch = true
					break
				}
			}

			if keywordMatch {
				break
			}

			// Check in article title and description
			titleLower := strings.ToLower(article.Title)
			descLower := strings.ToLower(article.Description)

			if strings.Contains(titleLower, alertKeywordLower) ||
			   strings.Contains(descLower, alertKeywordLower) {
				keywordMatch = true
				break
			}
		}

		if !keywordMatch {
			return false
		}
	}

	// Check source filter if specified
	if len(alert.SourceIDs) > 0 {
		sourceMatch := false
		for _, sourceID := range alert.SourceIDs {
			if uint(sourceID) == article.SourceID {
				sourceMatch = true
				break
			}
		}
		if !sourceMatch {
			return false
		}
	}

	return true
}