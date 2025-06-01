package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrrobotisreal/rss_today_api/internal/models"
)

// LoginWithTokenRequest represents the request body for login with Firebase ID token
type LoginWithTokenRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

func Login(app *models.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginWithTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		ctx := context.Background()

		// Verify the Firebase ID token
		token, err := app.FirebaseAuth.VerifyIDToken(ctx, req.IDToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid credentials: " + err.Error(),
			})
			return
		}

		// Find user in database
		var user models.User
		if err := app.DB.Where("firebase_uid = ?", token.UID).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "User not found in database",
			})
			return
		}

		// Generate custom token for the user (optional, since they already have an ID token)
		customToken, err := app.FirebaseAuth.CustomToken(ctx, token.UID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to generate authentication token: " + err.Error(),
			})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, models.AuthResponse{
			User:  user,
			Token: customToken,
		})
	}
}