package app

import (
	"extension-3/kafkaservice"
	"log"
	"sync"
)

type App struct {
	Mu             sync.Mutex
	UniqueIDCache  map[string]struct{}
	MinuteLogger   *log.Logger
	KafkaService   *kafkaservice.KafkaService
	ShutdownSignal chan struct{}
}
