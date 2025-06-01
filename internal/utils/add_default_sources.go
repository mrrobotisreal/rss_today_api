package utils

import (
	"log"

	"gorm.io/gorm"
)

func (app *models.App) AddDefaultSources() {
	defaultSources := []NewsSource{
		{Name: "BBC News", URL: "https://www.bbc.com/news", RSSURL: "http://feeds.bbci.co.uk/news/rss.xml"},
		{Name: "CNN", URL: "https://www.cnn.com", RSSURL: "http://rss.cnn.com/rss/edition.rss"},
		{Name: "Reuters", URL: "https://www.reuters.com", RSSURL: "https://www.reuters.com/rssFeed/topNews"},
		{Name: "Associated Press", URL: "https://apnews.com", RSSURL: "https://feeds.apnews.com/rss/topnews"},
	}

	for _, source := range defaultSources {
		var existing NewsSource
		result := app.DB.Where("rss_url = ?", source.RSSURL).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			app.DB.Create(&source)
			log.Printf("Added default source: %s", source.Name)
		}
	}
}
