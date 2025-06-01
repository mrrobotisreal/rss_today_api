package models

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron/v3"
	"google.golang.org/api/option"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	DB           *gorm.DB
	Router       *gin.Engine
	Parser       *gofeed.Parser
	Cron         *cron.Cron
	FirebaseAuth *auth.Client
	mu           sync.RWMutex
	lastRun      time.Time
}
