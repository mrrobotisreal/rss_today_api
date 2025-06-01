package handlers

import (
	"context"
	"net/http"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/mrrobotisreal/rss_today_api/internal/models"
)

func CreateUser(app *models.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		ctx := context.Background()

		// Create user in Firebase Auth
		params := (&auth.UserToCreate{}).
			Email(req.Email).
			Password(req.Password).
			EmailVerified(false).
			Disabled(false)

		if req.DisplayName != "" {
			params = params.DisplayName(req.DisplayName)
		}

		firebaseUser, err := app.FirebaseAuth.CreateUser(ctx, params)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Failed to create user: " + err.Error(),
			})
			return
		}

		// Create user in database
		user := models.User{
			FirebaseUID: firebaseUser.UID,
			Email:       req.Email,
			DisplayName: req.DisplayName,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := app.DB.Create(&user).Error; err != nil {
			// If database creation fails, clean up Firebase user
			app.FirebaseAuth.DeleteUser(ctx, firebaseUser.UID)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to save user to database: " + err.Error(),
			})
			return
		}

		// Generate custom token for the user
		token, err := app.FirebaseAuth.CustomToken(ctx, firebaseUser.UID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to generate authentication token: " + err.Error(),
			})
			return
		}

		// Return success response
		c.JSON(http.StatusCreated, models.AuthResponse{
			User:  user,
			Token: token,
		})
	}
}