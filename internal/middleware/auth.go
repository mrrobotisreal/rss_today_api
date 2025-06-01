package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mrrobotisreal/rss_today_api/internal/models"
	"gorm.io/gorm"
)

func AuthMiddleware(app *models.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		firebaseToken, err := app.FirebaseAuth.VerifyIDToken(context.Background(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Get or create user in our database
		var user models.User
		result := app.DB.Where("firebase_uid = ?", firebaseToken.UID).First(&user)
		if result.Error == gorm.ErrRecordNotFound {
			// Create new user
			user = models.User{
				FirebaseUID: firebaseToken.UID,
				Email:       firebaseToken.Claims["email"].(string),
				DisplayName: firebaseToken.Claims["name"].(string),
			}
			app.DB.Create(&user)
		}

		// Add user to context
		c.Set("user", user)
		c.Next()
	}
}
