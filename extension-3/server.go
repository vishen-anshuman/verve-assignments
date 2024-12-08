package main

import (
	"extension-3/app"
	"extension-3/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var appInstance *app.App

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

func initLogUniqueCount() {
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
			appInstance.Mu.Lock()
			count, _ := appInstance.RedisService.ReadFromCache(app.UNIQUE_COUNT)
			appInstance.RedisService.WriteToCache(app.UNIQUE_COUNT, "0", 0)
			appInstance.Mu.Unlock()
			appInstance.MinuteLogger.Printf("Unique requests in the last minute: %d", count)

		case <-appInstance.ShutdownSignal:
			log.Println("Shutdown signal received, stopping periodic logger.")
			return
		}
	}
}

func main() {
	app.InitApp()
	appInstance = app.GetAppConst()
	go initLogUniqueCount()
	http.HandleFunc("/api/verve/accept", handlers.AcceptHandler)

	port := 8080
	log.Printf("Starting server on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
