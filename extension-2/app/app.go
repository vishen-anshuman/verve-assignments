package app

import (
	"extension-2/logger"
	"extension-2/redisservice"
	"fmt"
	"log"
	"sync"
)

const (
	UNIQUE_COUNT         = "UNIQUE_COUNT"
	UNIQUE_ID_FORMAT     = "UNIQUE_ID_%s"
	PROCESSING_ID_FORMAT = "PROCESSING_ID_%s"
)

type App struct {
	Mu             sync.Mutex
	RedisService   *redisservice.RedisService
	MinuteLogger   *log.Logger
	ShutdownSignal chan struct{}
}

var appConst *App

func InitApp() {
	redisAddr := fmt.Sprintf("localhost:%d", 6379)
	appConst = &App{
		Mu:             sync.Mutex{},
		RedisService:   redisservice.InitRedisService(redisAddr),
		MinuteLogger:   logger.InitLogger(),
		ShutdownSignal: make(chan struct{}),
	}
}

func GetAppConst() *App {
	return appConst
}
