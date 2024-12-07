package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

func InitLogger() *log.Logger {
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
