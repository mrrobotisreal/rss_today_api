package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mrrobotisreal/rss_today_api/internal/models"
)

// VerifyTokenRequest represents the request body for token verification
type VerifyTokenRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

func VerifyToken(app *models.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VerifyTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		ctx := context.Background()

		// Verify the ID token
		token, err := app.FirebaseAuth.VerifyIDToken(ctx, req.IDToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid token: " + err.Error(),
			})
			return
		}

		// Find user in database
		var user models.User
		if err := app.DB.Where("firebase_uid = ?", token.UID).First(&user).Error; err != nil {
			// If user doesn't exist in database, create them
			firebaseUser, err := app.FirebaseAuth.GetUser(ctx, token.UID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{
					Error: "Failed to get user from Firebase: " + err.Error(),
				})
				return
			}

			user = models.User{
				FirebaseUID: firebaseUser.UID,
				Email:       firebaseUser.Email,
				DisplayName: firebaseUser.DisplayName,
			}

			if err := app.DB.Create(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{
					Error: "Failed to create user in database: " + err.Error(),
				})
				return
			}
		}

		// Return user information
		c.JSON(http.StatusOK, gin.H{
			"user": user,
			"uid":  token.UID,
		})
	}
}

// Helper function to extract token from Authorization header
func ExtractTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Bearer token format: "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}