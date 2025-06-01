package models

import "time"

type NewsSource struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`          // e.g. "BBC News"
	URL       string    `json:"url"`                           // e.g. "https://bbc.com"
	RSSURL    string    `json:"rss_url" gorm:"not null"`       // e.g. "http://feeds.bbci.co.uk/news/rss.xml"
	Active    bool      `json:"active" gorm:"default:true"`    // Whether to monitor this source
	CreatedAt time.Time `json:"created_at"`
}
