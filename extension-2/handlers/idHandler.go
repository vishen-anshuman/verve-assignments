package handlers

import (
	"extension-2/app"
	helperfunc "extension-2/helper"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func AcceptHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if status := helperfunc.BeforeProcessing(r); status == 400 || status == 409 {
			if status == 409 {
				http.Error(w, "Requested Id being Processed", http.StatusConflict)
				return
			}
			http.Error(w, "Invalid Requested Id", http.StatusBadRequest)
			return
		}
		query := r.URL.Query()
		idParam := query.Get("id")
		endpoint := query.Get("endpoint")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "failed", http.StatusBadRequest)
			return
		}

		log.Printf("Received id: %d, endpoint: %s", id, endpoint)
		if endpoint != "" {
			go fireEndpointRequest(endpoint, len(app.UniqueIDCache), app.MinuteLogger)
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("ok"))
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
		defer func() {
			helperfunc.AfterProcessing(idParam)
		}()
	}
}

func fireEndpointRequest(endpoint string, count int, logger *log.Logger) {
	url := fmt.Sprintf("%s?count=%d", endpoint, count)
	resp, err := http.Get(url)
	if err != nil {
		logger.Printf("Error making GET request to %s: %v", url, err)
		log.Printf("Error making GET request to %s: %v", url, err)
		return
	}
	defer resp.Body.Close()
	logger.Printf("GET request to %s returned status code: %d", url, resp.StatusCode)
	log.Printf("GET request to %s returned status code: %d", url, resp.StatusCode)
}
