package services

import (
	"crypto/sha256"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mrrobotisreal/rss_today_api/internal/models"
)

func FetchRSSFeed(app *models.App, source models.NewsSource) ([]models.Article, error) {
	log.Printf("Fetching RSS feed for %s", source.Name)

	feed, err := app.Parser.ParseURL(source.RSSURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing RSS feed for %s: %v", source.Name, err)
	}

	var articles []models.Article
	for _, item := range feed.Items {
		if item.Title == "" || item.Link == "" {
			continue
		}

		// Create content hash for duplicate detection
		contentData := item.Title + item.Link + item.Description
		hash := sha256.Sum256([]byte(contentData))
		contentHash := fmt.Sprintf("%x", hash)

		// Parse publication date
		var pubDate time.Time
		if item.PublishedParsed != nil {
			pubDate = *item.PublishedParsed
		} else {
			pubDate = time.Now()
		}

		// Extract basic keywords from title and description
		keywords := extractKeywords(item.Title + " " + item.Description)

		article := models.Article{
			SourceID:    source.ID,
			Title:       item.Title,
			Description: item.Description,
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

func extractKeywords(text string) []string {
	// Simple keyword extraction - split by spaces and clean up
	words := strings.Fields(strings.ToLower(text))
	var keywords []string

	for _, word := range words {
		// Remove punctuation and keep words longer than 3 characters
		cleaned := strings.Trim(word, ".,!?:;\"'()[]{}*-_+=<>/\\|")
		if len(cleaned) > 3 && !isStopWord(cleaned) {
			keywords = append(keywords, cleaned)
		}
	}

	// Remove duplicates and limit to 10 keywords
	uniqueKeywords := removeDuplicates(keywords)
	if len(uniqueKeywords) > 10 {
		uniqueKeywords = uniqueKeywords[:10]
	}

	return uniqueKeywords
}

func isStopWord(word string) bool {
	stopWords := map[string]bool{
		"the": true, "and": true, "for": true, "are": true, "but": true,
		"not": true, "you": true, "all": true, "can": true, "had": true,
		"her": true, "was": true, "one": true, "our": true, "out": true,
		"day": true, "get": true, "him": true, "his": true,
		"how": true, "its": true, "may": true, "new": true, "now": true,
		"old": true, "see": true, "two": true, "who": true, "boy": true,
		"did": true, "let": true, "put": true, "say": true,
		"she": true, "too": true, "use": true, "with": true, "that": true,
		"this": true, "will": true, "have": true, "from": true, "they": true,
		"know": true, "want": true, "been": true, "good": true, "much": true,
		"some": true, "time": true, "very": true, "when": true, "come": true,
		"here": true, "just": true, "like": true, "long": true, "make": true,
		"many": true, "over": true, "such": true, "take": true, "than": true,
		"them": true, "well": true, "were": true, "has": true,
	}
	return stopWords[word]
}

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	return result
}