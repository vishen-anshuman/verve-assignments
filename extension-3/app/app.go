package app

import (
	"extension-3/kafkaservice"
	"extension-3/logger"
	"extension-3/redisservice"
	"fmt"
	"log"
	"sync"
)

type App struct {
	Mu             sync.Mutex
	RedisService   *redisservice.RedisService
	MinuteLogger   *log.Logger
	KafkaService   *kafkaservice.KafkaService
	ShutdownSignal chan struct{}
}

const (
	UNIQUE_COUNT         = "UNIQUE_COUNT"
	UNIQUE_ID_FORMAT     = "UNIQUE_ID_%s"
	PROCESSING_ID_FORMAT = "PROCESSING_ID_%s"
)

var appConst *App

func InitApp() {
	brokers := []string{"localhost:9092"}
	topic := "verve-streaming"
	redisAddr := fmt.Sprintf("localhost:%d", 6379)
	appConst = &App{
		Mu:             sync.Mutex{},
		RedisService:   redisservice.InitRedisService(redisAddr),
		MinuteLogger:   logger.InitLogger(),
		KafkaService:   kafkaservice.InitKafkaService(brokers, topic),
		ShutdownSignal: make(chan struct{}),
	}
}

func GetAppConst() *App {
	return appConst
}
