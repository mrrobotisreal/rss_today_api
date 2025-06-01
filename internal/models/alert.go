package models

import (
	"time"

	"github.com/lib/pq"
)

type UserAlert struct {
	ID                  uint           `json:"id" gorm:"primaryKey"`
	UserID              uint           `json:"user_id" gorm:"not null"`                        // Which user
	Keywords            pq.StringArray `json:"keywords" gorm:"type:text[]"`                    // Keywords to watch for ["ukraine", "war"]
	SourceIDs           pq.Int64Array  `json:"source_ids" gorm:"type:integer[]"`               // Which sources to monitor (empty = all)
	NotificationMethods pq.StringArray `json:"notification_methods" gorm:"type:text[]"`       // ["email", "push", "sms"]
	Active              bool           `json:"active" gorm:"default:true"`                     // Whether alert is enabled
	CreatedAt           time.Time      `json:"created_at"`
}
