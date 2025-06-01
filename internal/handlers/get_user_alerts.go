package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrrobotisreal/rss_today_api/internal/models"
)

func (app *models.App) GetUserAlerts(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(models.User)

	var alerts []models.UserAlert
	if err := app.DB.Where("user_id = ?", currentUser.ID).Find(&alerts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alerts)
}
