package models

import "time"

type NotificationSent struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`    // Which user was notified
	ArticleID uint      `json:"article_id" gorm:"not null"` // Which article triggered notification
	AlertID   uint      `json:"alert_id" gorm:"not null"`   // Which alert rule matched
	Method    string    `json:"method" gorm:"not null"`     // "email", "push", "sms"
	SentAt    time.Time `json:"sent_at"`                    // When notification was sent
}
