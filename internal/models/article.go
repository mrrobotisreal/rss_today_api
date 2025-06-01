package models

import (
	"time"

	"github.com/lib/pq"
)

type Article struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	SourceID    uint           `json:"source_id" gorm:"not null"`                    // Which news source
	Title       string         `json:"title" gorm:"not null"`                       // Article headline
	Description string         `json:"description"`                                 // Article summary
	Link        string         `json:"link" gorm:"unique;not null"`                 // Original article URL
	PubDate     time.Time      `json:"pub_date"`                                    // When published
	ContentHash string         `json:"content_hash" gorm:"unique"`                  // Hash to detect duplicates
	Keywords    pq.StringArray `json:"keywords" gorm:"type:text[]"`                 // Extracted keywords
	CreatedAt   time.Time      `json:"created_at"`                                  // When we found it
	Source      NewsSource     `json:"source,omitempty" gorm:"foreignKey:SourceID"` // Join with source
}
