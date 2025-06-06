package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrrobotisreal/rss_today_api/internal/models"
)

func GetSources(app *models.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sources []models.NewsSource
		if err := app.DB.Find(&sources).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, sources)
	}
}
