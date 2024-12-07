package handlers

import (
	"extension-3/app"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func AcceptHandler(w http.ResponseWriter, r *http.Request) {
	appInstance := app.GetAppConst()
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
	uniqueCount, _ := appInstance.RedisService.ReadFromCache("UNIQUE_COUNT")
	endpoint := query.Get("endpoint")
	log.Printf("Received id: %d, endpoint: %s", id, endpoint)
	if endpoint != "" {
		go fireEndpointRequest(endpoint, uniqueCount, appInstance.MinuteLogger)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("ok"))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
	appInstance.Mu.Lock()
	idParamKey := fmt.Sprintf("UNIQUE_COUNT_%s", idParam)
	idCount, _ := appInstance.RedisService.ReadFromCache(idParamKey)
	if idCount == "" {
		appInstance.RedisService.WriteToCache(idParamKey, "1")
		countInt, _ := strconv.Atoi(uniqueCount)
		appInstance.RedisService.WriteToCache("UNIQUE_COUNT", string(countInt+1))
	}
	appInstance.Mu.Unlock()
}

func fireEndpointRequest(endpoint, count string, logger *log.Logger) {
	url := fmt.Sprintf("%s?count=%s", endpoint, count)
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
