package services

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func (app *models.App) MonitorAllSources() error {
	app.mu.Lock()
	app.lastRun = time.Now()
	app.mu.Unlock()

	log.Println("🔍 STARTING RSS MONITORING CYCLE...")

	// Step 1: Get all active news sources from database
	var sources []NewsSource
	if err := app.DB.Where("active = ?", true).Find(&sources).Error; err != nil {
		return fmt.Errorf("error fetching sources: %v", err)
	}

	log.Printf("Monitoring %d RSS sources", len(sources))

	// Step 2: Process all sources concurrently using goroutines
	var allNewArticles []Article
	var wg sync.WaitGroup
	articlesChan := make(chan []Article, len(sources))

	for _, source := range sources {
		wg.Add(1)
		go func(src NewsSource) {
			defer wg.Done()

			// Fetch RSS feed
			articles, err := app.FetchRSSFeed(src)
			if err != nil {
				log.Printf("Error fetching RSS for %s: %v", src.Name, err)
				return
			}

			// Save new articles to database
			newArticles, err := app.SaveNewArticles(articles)
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
		if err := app.CheckAlertsForNewArticles(allNewArticles); err != nil {
			log.Printf("Error checking alerts: %v", err)
		}
	} else {
		log.Println("📰 No new articles found this cycle")
	}

	log.Println("✅ RSS monitoring cycle completed")
	return nil
}
