package app

import (
	"extension-2/logger"
	"log"
	"sync"
)

type App struct {
	Mu             sync.Mutex
	UniqueIDCache  map[string]struct{}
	MinuteLogger   *log.Logger
	ShutdownSignal chan struct{}
}

var appConst *App

func InitApp() {
	appConst = &App{
		Mu:             sync.Mutex{},
		UniqueIDCache:  make(map[string]struct{}),
		MinuteLogger:   logger.InitLogger(),
		ShutdownSignal: make(chan struct{}),
	}
}

func GetAppConst() *App {
	return appConst
}
