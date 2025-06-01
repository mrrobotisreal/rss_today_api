package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrrobotisreal/rss_today_api/internal/models"
)

func (app *models.App) triggerMonitoring(c *gin.Context) {
	go func() {
		if err := app.MonitorAllSources(); err != nil {
			log.Printf("Error in manual monitoring: %v", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Monitoring triggered successfully"})
}
