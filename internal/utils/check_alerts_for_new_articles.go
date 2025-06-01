package utils

import (
	"fmt"
	"log"
)

func (app *models.App) CheckAlertsForNewArticles(newArticles []Article) error {
	if len(newArticles) == 0 {
		return nil
	}

	// Get all active alerts with user info
	var alerts []UserAlert
	if err := app.DB.Where("active = ?", true).Find(&alerts).Error; err != nil {
		return fmt.Errorf("error fetching alerts: %v", err)
	}

	log.Printf("Checking %d new articles against %d active alerts", len(newArticles), len(alerts))

	for _, alert := range alerts {
		// Get user info
		var user User
		if err := app.DB.First(&user, alert.UserID).Error; err != nil {
			log.Printf("Error finding user %d: %v", alert.UserID, err)
			continue
		}

		for _, article := range newArticles {
			// Check if article matches alert keywords
			if !app.ArticleMatchesAlert(article, alert) {
				continue
			}

			// Check if user wants notifications from this specific source
			if len(alert.SourceIDs) > 0 {
				sourceMatched := false
				for _, sourceID := range alert.SourceIDs {
					if int64(article.SourceID) == sourceID {
						sourceMatched = true
						break
					}
				}
				if !sourceMatched {
					continue
				}
			}

			// Send notification
			if err := app.SendNotification(user, alert, article); err != nil {
				log.Printf("Error sending notification: %v", err)
			}
		}
	}

	return nil
}
