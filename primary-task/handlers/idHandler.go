package handlers

import (
	"fmt"
	"log"
	"net/http"
	"primary-task/app"
	"strconv"
)

func AcceptHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		idParam := query.Get("id")
		if idParam == "" {
			http.Error(w, "failed", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "failed", http.StatusBadRequest)
			return
		}

		endpoint := query.Get("endpoint")
		log.Printf("Received id: %d, endpoint: %s", id, endpoint)
		if endpoint != "" {
			go fireEndpointRequest(endpoint, len(app.UniqueIDCache), app.MinuteLogger)
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("ok"))
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
		app.Mu.Lock()
		app.UniqueIDCache[idParam] = struct{}{}
		app.Mu.Unlock()
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
