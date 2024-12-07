package main

import (
	"extension-3/app"
	"extension-3/handlers"
	"extension-3/kafkaservice"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func setupLogger() *log.Logger {
	timestamp := time.Now().Format("20060102_150405")
	logFileName := fmt.Sprintf("/app/logs/unique_requests_%s.log", timestamp)
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	logger := log.New(logFile, "", log.Ldate|log.Ltime)
	log.Printf("Minute logger initialized: %s", logFileName)
	return logger
}

func initLogUniqueCount(app *app.App) {
	// I wish to have a new ticker independent of starting time,
	//suppose if it starts at 11:22:23, then new should come at 11:23:00, not 11:23:23,
	//and thereafter it should be 11:24:00, 11:25:00 etc.
	now := time.Now()
	nextFullMinute := now.Truncate(time.Minute).Add(time.Minute)
	durationUntilNextMinute := nextFullMinute.Sub(now)
	time.Sleep(durationUntilNextMinute)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			app.Mu.Lock()
			count := len(app.UniqueIDCache)
			app.UniqueIDCache = make(map[string]struct{})
			app.KafkaService.WriteLog(fmt.Sprintf("Unique requests in the last minute: %d", count))
			app.Mu.Unlock()

		case <-app.ShutdownSignal:
			log.Println("Shutdown signal received, stopping periodic logger.")
			return
		}
	}
}

func main() {
	brokers := []string{"localhost:9092"}
	topic := "example-topic"
	appConfig := &app.App{
		UniqueIDCache:  make(map[string]struct{}),
		KafkaService:   kafkaservice.InitKafkaService(brokers, topic),
		MinuteLogger:   setupLogger(),
		ShutdownSignal: make(chan struct{}),
	}
	go initLogUniqueCount(appConfig)
	http.HandleFunc("/api/verve/accept", handlers.AcceptHandler(appConfig))

	port := 8080
	log.Printf("Starting server on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
