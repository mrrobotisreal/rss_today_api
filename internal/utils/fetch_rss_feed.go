package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/mmcdole/gofeed/rss"
)

func (app *models.App) FetchRSSFeed(source NewsSource) ([]Article, error) {
	log.Printf("Fetching RSS for %s...", source.Name)

	feed, err := app.Parser.ParseURL(source.RSSURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing RSS for %s: %v", source.Name, err)
	}

	var articles []Article
	for _, item := range feed.Items {
		if item.Title == "" || item.Link == "" {
			continue
		}

		// Parse publication date
		pubDate := time.Now()
		if item.PublishedParsed != nil {
			pubDate = *item.PublishedParsed
		}

		// Get description
		description := ""
		if item.Description != "" {
			description = item.Description
		} else if item.Content != "" {
			description = item.Content
		}

		// Extract keywords and generate hash
		keywords := app.ExtractKeywords(item.Title, description)
		contentHash := app.GenerateContentHash(item.Title, item.Link)

		article := Article{
			SourceID:    source.ID,
			Title:       item.Title,
			Description: description,
			Link:        item.Link,
			PubDate:     pubDate,
			ContentHash: contentHash,
			Keywords:    keywords,
		}

		articles = append(articles, article)
	}

	log.Printf("Parsed %d articles from %s", len(articles), source.Name)
	return articles, nil
}
