package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrrobotisreal/rss_today_api/internal/models"
	"github.com/mrrobotisreal/rss_today_api/internal/services"
)

func TriggerMonitoring(app *models.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		go func() {
			if err := services.MonitorAllSources(app); err != nil {
				log.Printf("Error in manual monitoring: %v", err)
			}
		}()

		c.JSON(http.StatusOK, gin.H{"message": "Monitoring triggered successfully"})
	}
}
