package models

import (
	"sync"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type App struct {
	DB           *gorm.DB
	Router       *gin.Engine
	Parser       *gofeed.Parser
	Cron         *cron.Cron
	FirebaseAuth *auth.Client
	Mu           sync.RWMutex
	LastRun      time.Time
}
