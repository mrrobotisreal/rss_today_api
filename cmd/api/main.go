package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"github.com/mrrobotisreal/rss_today_api/internal/db"
	"github.com/mrrobotisreal/rss_today_api/internal/firebase"
	"github.com/mrrobotisreal/rss_today_api/internal/handlers"
	"github.com/mrrobotisreal/rss_today_api/internal/middleware"
	"github.com/mrrobotisreal/rss_today_api/internal/models"
	"github.com/mrrobotisreal/rss_today_api/internal/services"
	"github.com/robfig/cron/v3"
)

// Setup all routes
func setupRoutes(app *models.App) {
	app.Router = gin.Default()

	// CORS
	app.Router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Public routes
	app.Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Protected routes (require Firebase authentication)
	api := app.Router.Group("/api")
	api.Use(middleware.AuthMiddleware(app))
	{
		api.GET("/articles", handlers.GetArticles(app))
		api.GET("/sources", handlers.GetSources(app))
		api.POST("/alerts", handlers.CreateAlert(app))
		api.GET("/alerts", handlers.GetUserAlerts(app))
		api.POST("/monitor/trigger", handlers.TriggerMonitoring(app))
	}
}

// Start the cron scheduler
func startScheduler(app *models.App) {
	app.Cron = cron.New()

	// Run every 10 minutes: "*/10 * * * *"
	app.Cron.AddFunc("*/10 * * * *", func() {
		log.Println("⏰ CRON JOB TRIGGERED - Running RSS monitoring...")
		if err := services.MonitorAllSources(app); err != nil {
			log.Printf("Error in scheduled monitoring: %v", err)
		}
	})

	app.Cron.Start()
	log.Println("📡 Cron scheduler started - RSS monitoring every 10 minutes")
}

func main() {
	app := &models.App{
		Parser: gofeed.NewParser(),
	}

	// Initialize Firebase
	if err := firebase.InitFirebase(app); err != nil {
		log.Fatal("Failed to initialize Firebase:", err)
	}

	// Initialize database
	if err := db.InitDatabase(app); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Setup routes
	setupRoutes(app)

	// Start cron scheduler
	startScheduler(app)

	// Run initial monitoring after 10 seconds
	go func() {
		time.Sleep(10 * time.Second)
		log.Println("Running initial RSS monitoring...")
		if err := services.MonitorAllSources(app); err != nil {
			log.Printf("Error in initial monitoring: %v", err)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	log.Printf("🚀 RSS Monitor API running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, app.Router))
}
