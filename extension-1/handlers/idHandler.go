package handlers

import (
	"bytes"
	"encoding/json"
	"extension-1/app"
	"log"
	"net/http"

	"strconv"
	"time"
)

func AcceptHandler(w http.ResponseWriter, r *http.Request) {
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
	appConst := app.GetAppConst()
	endpoint := query.Get("endpoint")
	log.Printf("Received id: %d, endpoint: %s", id, endpoint)
	if endpoint != "" {
		go fireEndpointRequest(endpoint, len(appConst.UniqueIDCache), appConst.MinuteLogger)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("ok"))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
	appConst.Mu.Lock()
	appConst.UniqueIDCache[idParam] = struct{}{}
	appConst.Mu.Unlock()
}

func fireEndpointRequest(endpoint string, count int, logger *log.Logger) {
	payload := map[string]interface{}{
		"count":     count,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	requestBody, err := json.Marshal(payload)
	if err != nil {
		logger.Printf("Error marshalling JSON payload: %v", err)
		log.Printf("Error marshalling JSON payload: %v", err)
		return
	}
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		logger.Printf("Error making POST request to %s: %v", endpoint, err)
		log.Printf("Error making POST request to %s: %v", endpoint, err)
		return
	}
	defer resp.Body.Close()
	logger.Printf("POST request to %s returned status code: %d", endpoint, resp.StatusCode)
	log.Printf("POST request to %s returned status code: %d", endpoint, resp.StatusCode)
}
