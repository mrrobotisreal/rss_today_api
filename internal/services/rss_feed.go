package services

import (
	"crypto/sha256"
	"fmt"
	"html"
	"log"
	"net/url"
	"regexp"
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

	// Initialize Google News decoder for this source if needed
	var decoder *GoogleNewsDecoder
	if strings.Contains(source.Name, "Google News") || strings.Contains(source.RSSURL, "news.google.com") {
		decoder = NewGoogleNewsDecoder()
		log.Printf("Initialized Google News decoder for %s", source.Name)
	}

	var articles []models.Article
	for _, item := range feed.Items {
		if item.Title == "" || item.Link == "" {
			continue
		}

		// Clean and decode title
		cleanTitle := CleanContent(item.Title)

		// Clean and decode description
		cleanDescription := CleanContent(item.Description)

		// Process the link - decode Google News URLs if needed
		cleanLink := processLink(item.Link, source.Name, decoder)

		// Create content hash for duplicate detection
		contentData := cleanTitle + cleanLink + cleanDescription
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
		keywords := extractKeywords(cleanTitle + " " + cleanDescription)

		article := models.Article{
			SourceID:    source.ID,
			Title:       cleanTitle,
			Description: cleanDescription,
			Link:        cleanLink,
			PubDate:     pubDate,
			ContentHash: contentHash,
			Keywords:    keywords,
		}

		articles = append(articles, article)
	}

	log.Printf("Parsed %d articles from %s", len(articles), source.Name)
	return articles, nil
}

// CleanContent removes HTML tags, decodes HTML entities, and cleans up text
func CleanContent(content string) string {
	if content == "" {
		return ""
	}

	// Decode HTML entities (like &amp;, &lt;, &gt;, &quot;, etc.)
	content = html.UnescapeString(content)

	// Remove HTML tags
	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	content = htmlTagRegex.ReplaceAllString(content, "")

	// Clean up extra whitespace
	content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")
	content = strings.TrimSpace(content)

	// Remove source name from end of Google News titles (like "... - BBC News")
	sourcePattern := regexp.MustCompile(` - [^-]+$`)
	content = sourcePattern.ReplaceAllString(content, "")

	return content
}

// processLink handles Google News URL decoding and other link processing
func processLink(link, sourceName string, decoder *GoogleNewsDecoder) string {
	if link == "" {
		return ""
	}

	// Check if this is a Google News redirect URL and we have a decoder
	if decoder != nil && IsGoogleNewsURL(link) {
		if decodedURL, err := decoder.DecodeGoogleNewsURL(link); err == nil && decodedURL != link {
			log.Printf("Successfully decoded Google News URL for %s: %s -> %s", sourceName, link, decodedURL)
			return decodedURL
		} else {
			log.Printf("Could not decode Google News URL for %s: %s", sourceName, link)
		}
	}

	// Clean up other URL issues
	parsedURL, err := url.Parse(link)
	if err != nil {
		log.Printf("Error parsing URL %s: %v", link, err)
		return link
	}

	// Remove tracking parameters from URLs
	cleanQuery := url.Values{}
	for key, values := range parsedURL.Query() {
		// Keep important parameters, remove tracking ones
		if !isTrackingParameter(key) {
			cleanQuery[key] = values
		}
	}
	parsedURL.RawQuery = cleanQuery.Encode()

	return parsedURL.String()
}

// isTrackingParameter identifies common tracking parameters to remove
func isTrackingParameter(param string) bool {
	trackingParams := map[string]bool{
		"utm_source":   true,
		"utm_medium":   true,
		"utm_campaign": true,
		"utm_term":     true,
		"utm_content":  true,
		"fbclid":       true,
		"gclid":        true,
		"ref":          true,
		"source":       true,
		"oc":           true, // Google News specific parameter
	}
	return trackingParams[strings.ToLower(param)]
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