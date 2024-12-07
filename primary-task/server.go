package main

import (
	"fmt"
	"log"
	"net/http"
	"primary-task/app"
	"primary-task/handlers"
	"time"
)

var appInstance *app.App

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
			count := len(appInstance.UniqueIDCache)
			appInstance.UniqueIDCache = make(map[string]struct{})
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
