package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrrobotisreal/rss_today_api/internal/models"
	"github.com/mrrobotisreal/rss_today_api/internal/services"
)

// TestGoogleNewsResponse represents the response for testing Google News processing
type TestGoogleNewsResponse struct {
	Source           string    `json:"source"`
	TotalArticles    int       `json:"total_articles"`
	ProcessedArticles []ProcessedArticle `json:"processed_articles"`
	Errors           []string  `json:"errors,omitempty"`
}

// ProcessedArticle represents a processed article with before/after comparison
type ProcessedArticle struct {
	OriginalTitle       string `json:"original_title"`
	CleanedTitle        string `json:"cleaned_title"`
	OriginalDescription string `json:"original_description"`
	CleanedDescription  string `json:"cleaned_description"`
	OriginalLink        string `json:"original_link"`
	ProcessedLink       string `json:"processed_link"`
	IsGoogleNewsURL     bool   `json:"is_google_news_url"`
	URLDecoded          bool   `json:"url_decoded"`
}

// TestGoogleNews tests Google News RSS feed processing with enhanced decoding
func TestGoogleNews(app *models.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Google News source from the database
		var googleNewsSource models.NewsSource
		result := app.DB.Where("name LIKE ? OR rss_url LIKE ?", "%Google News%", "%news.google.com%").First(&googleNewsSource)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Google News source not found in database",
				"suggestion": "Make sure you have a Google News source configured",
			})
			return
		}

		// Parse the RSS feed
		feed, err := app.Parser.ParseURL(googleNewsSource.RSSURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to parse Google News RSS feed",
				"details": err.Error(),
			})
			return
		}

		// Initialize decoder
		decoder := services.NewGoogleNewsDecoder()

		var processedArticles []ProcessedArticle
		var errors []string

		// Process up to 10 articles for testing
		maxArticles := 10
		if len(feed.Items) < maxArticles {
			maxArticles = len(feed.Items)
		}

		for i := 0; i < maxArticles; i++ {
			item := feed.Items[i]

			// Process the article
			isGoogleNewsURL := services.IsGoogleNewsURL(item.Link)
			processedLink := item.Link
			urlDecoded := false

			if isGoogleNewsURL {
				if decodedURL, err := decoder.DecodeGoogleNewsURL(item.Link); err == nil && decodedURL != item.Link {
					processedLink = decodedURL
					urlDecoded = true
				}
			}

			// Clean content using the public function from services
			cleanedTitle := services.CleanContent(item.Title)
			cleanedDescription := services.CleanContent(item.Description)

			processedArticle := ProcessedArticle{
				OriginalTitle:       item.Title,
				CleanedTitle:        cleanedTitle,
				OriginalDescription: item.Description,
				CleanedDescription:  cleanedDescription,
				OriginalLink:        item.Link,
				ProcessedLink:       processedLink,
				IsGoogleNewsURL:     isGoogleNewsURL,
				URLDecoded:          urlDecoded,
			}

			processedArticles = append(processedArticles, processedArticle)
		}

		response := TestGoogleNewsResponse{
			Source:            googleNewsSource.Name,
			TotalArticles:     len(feed.Items),
			ProcessedArticles: processedArticles,
			Errors:           errors,
		}

		c.JSON(http.StatusOK, response)
	}
}