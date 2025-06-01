package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/mrrobotisreal/rss_today_api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func GetArticles(app *models.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		keywords := c.Query("keywords")
		sourceID := c.Query("source_id")
		limitStr := c.DefaultQuery("limit", "50")

		limit, _ := strconv.Atoi(limitStr)

		query := app.DB.Model(&models.Article{}).Preload("Source")

		if keywords != "" {
			keywordList := strings.Split(keywords, ",")
			for i, keyword := range keywordList {
				keywordList[i] = strings.TrimSpace(keyword)
			}
			query = query.Where("keywords && ?", pq.Array(keywordList))
		}

		if sourceID != "" {
			query = query.Where("source_id = ?", sourceID)
		}

		var articles []models.Article
		if err := query.Order("pub_date DESC").Limit(limit).Find(&articles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, articles)
	}
}
