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

var appConst *App

func InitApp() {
	brokers := []string{"localhost:9092"}
	topic := "example-topic"
	redisAddr := fmt.Sprintf("localhost:%d", 6783)
	password := "hello123"
	appConst = &App{
		Mu:             sync.Mutex{},
		RedisService:   redisservice.InitRedisService(redisAddr, password, 0),
		MinuteLogger:   logger.InitLogger(),
		KafkaService:   kafkaservice.InitKafkaService(brokers, topic),
		ShutdownSignal: make(chan struct{}),
	}
}

func GetAppConst() *App {
	return appConst
}
