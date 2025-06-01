package services

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mrrobotisreal/rss_today_api/internal/models"
)

func MonitorAllSources(app *models.App) error {
	app.Mu.Lock()
	app.LastRun = time.Now()
	app.Mu.Unlock()

	log.Println("🔍 STARTING RSS MONITORING CYCLE...")

	// Step 1: Get all active news sources from database
	var sources []models.NewsSource
	if err := app.DB.Where("active = ?", true).Find(&sources).Error; err != nil {
		return fmt.Errorf("error fetching sources: %v", err)
	}

	log.Printf("Monitoring %d RSS sources", len(sources))

	// Step 2: Process all sources concurrently using goroutines
	var allNewArticles []models.Article
	var wg sync.WaitGroup
	articlesChan := make(chan []models.Article, len(sources))

	for _, source := range sources {
		wg.Add(1)
		go func(src models.NewsSource) {
			defer wg.Done()

			// Fetch RSS feed
			articles, err := FetchRSSFeed(app, src)
			if err != nil {
				log.Printf("Error fetching RSS for %s: %v", src.Name, err)
				return
			}

			// Save new articles to database
			newArticles, err := SaveNewArticles(app, articles)
			if err != nil {
				log.Printf("Error saving articles for %s: %v", src.Name, err)
				return
			}

			// Send new articles to channel if any found
			if len(newArticles) > 0 {
				articlesChan <- newArticles
			}

			// Be respectful - small delay between requests
			time.Sleep(1 * time.Second)
		}(source)
	}

	// Step 3: Collect all new articles
	go func() {
		wg.Wait()
		close(articlesChan)
	}()

	for articles := range articlesChan {
		allNewArticles = append(allNewArticles, articles...)
	}

	// Step 4: Check alerts and send notifications
	if len(allNewArticles) > 0 {
		log.Printf("📊 Found %d new articles total", len(allNewArticles))
		if err := CheckAlertsForNewArticles(app, allNewArticles); err != nil {
			log.Printf("Error checking alerts: %v", err)
		}
	} else {
		log.Println("📰 No new articles found this cycle")
	}

	log.Println("✅ RSS monitoring cycle completed")
	return nil
}
