package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mrrobotisreal/rss_today_api/internal/models"
)

// AuthInfoResponse represents Firebase configuration for client-side authentication
type AuthInfoResponse struct {
	ProjectID string `json:"project_id"`
	Message   string `json:"message"`
}

func GetAuthInfo(app *models.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := os.Getenv("FIREBASE_PROJECT_ID")
		if projectID == "" {
			projectID = "your-firebase-project-id" // fallback
		}

		response := AuthInfoResponse{
			ProjectID: projectID,
			Message:   "Use Firebase Client SDK for authentication. After authentication, send the ID token to /auth/login or /auth/verify endpoints.",
		}

		c.JSON(http.StatusOK, response)
	}
}

// SimpleLoginRequest for traditional email/password (returns instructions)
type SimpleLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// SimpleLogin provides instructions for client-side authentication
func SimpleLogin(app *models.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SimpleLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		// Return instructions for client-side authentication
		c.JSON(http.StatusOK, gin.H{
			"message": "Please use Firebase Client SDK to authenticate with these credentials, then send the ID token to /auth/login endpoint",
			"email":   req.Email,
			"next_step": "Use Firebase signInWithEmailAndPassword() and send the resulting ID token to /auth/login",
		})
	}
}