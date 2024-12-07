package app

import (
	"log"
	"sync"
)

type App struct {
	Mu             sync.Mutex
	UniqueIDCache  map[string]struct{}
	MinuteLogger   *log.Logger
	ShutdownSignal chan struct{}
}
