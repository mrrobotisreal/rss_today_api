package utils

import (
	"log"
	"time"

	"gorm.io/gorm"
)

func (app *models.App) SendNotification(user User, alert UserAlert, article Article) error {
	// Check if notification already sent to avoid spam
	var existing NotificationSent
	result := app.DB.Where("user_id = ? AND article_id = ? AND alert_id = ?",
		user.ID, article.ID, alert.ID).First(&existing)

	if result.Error != gorm.ErrRecordNotFound {
		return nil // Already sent
	}

	log.Printf("🚨 SENDING ALERT to %s (%s)", user.Email, user.DisplayName)
	log.Printf("📰 Article: %s", article.Title)
	log.Printf("🔗 URL: %s", article.Link)

	// Send notifications based on user's preferred methods
	for _, method := range alert.NotificationMethods {
		switch method {
		case "email":
			// TODO: Send email notification
			log.Printf("📧 Would send email to %s", user.Email)
		case "push":
			// TODO: Send push notification
			log.Printf("🔔 Would send push notification to %s", user.Email)
		case "sms":
			// TODO: Send SMS notification
			log.Printf("📱 Would send SMS to %s", user.Email)
		}

		// Record that notification was sent
		notification := NotificationSent{
			UserID:    user.ID,
			ArticleID: article.ID,
			AlertID:   alert.ID,
			Method:    method,
			SentAt:    time.Now(),
		}
		app.DB.Create(&notification)
	}

	return nil
}
