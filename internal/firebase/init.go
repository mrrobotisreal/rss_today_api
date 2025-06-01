package firebase

import (
	"context"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/auth"
	"github.com/mrrobotisreal/rss_today_api/internal/models"
	"google.golang.org/api/option"
)

func (app *models.App) initFirebase() error {
	ctx := context.Background()

	// Initialize Firebase with service account key
	serviceAccountKey := os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY")
	if serviceAccountKey == "" {
		return fmt.Errorf("FIREBASE_SERVICE_ACCOUNT_KEY environment variable is required")
	}

	opt := option.WithCredentialsFile(serviceAccountKey)
	firebaseApp, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing Firebase app: %v", err)
	}

	authClient, err := firebaseApp.Auth(ctx)
	if err != nil {
		return fmt.Errorf("error getting Auth client: %v", err)
	}

	app.FirebaseAuth = authClient
	log.Println("Firebase Auth initialized successfully")
	return nil
}
