package handlers

import (
	"extension-2/app"
	helperfunc "extension-2/helper"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func AcceptHandler(w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("ok"))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
	addUniqueCountToCache(idParam)
	if endpoint != "" {
		go fireEndpointRequest(endpoint)
	}
	defer func() {
		helperfunc.AfterProcessing(idParam)
	}()
}

func addUniqueCountToCache(idParam string) {
	appInstance := app.GetAppConst()
	appInstance.Mu.Lock()
	idParamKey := fmt.Sprintf(app.UNIQUE_ID_FORMAT, idParam)
	idCount, _ := appInstance.RedisService.ReadFromCache(idParamKey)
	uniqueCount, _ := appInstance.RedisService.ReadFromCache(app.UNIQUE_COUNT)
	if idCount == "" {
		appInstance.RedisService.WriteToCache(idParamKey, "exists", 0)
		countInt, _ := strconv.Atoi(uniqueCount)
		appInstance.RedisService.WriteToCache(app.UNIQUE_COUNT, string(countInt+1), 0)
	}
	appInstance.Mu.Unlock()
}

func fireEndpointRequest(endpoint string) {
	appInstance := app.GetAppConst()
	count, _ := appInstance.RedisService.ReadFromCache(app.UNIQUE_COUNT)
	url := fmt.Sprintf("%s?count=%d", endpoint, count)
	resp, err := http.Get(url)
	if err != nil {
		appInstance.MinuteLogger.Printf("Error making GET request to %s: %v", url, err)
		log.Printf("Error making GET request to %s: %v", url, err)
		return
	}
	defer resp.Body.Close()
	appInstance.MinuteLogger.Printf("GET request to %s returned status code: %d", url, resp.StatusCode)
	log.Printf("GET request to %s returned status code: %d", url, resp.StatusCode)
}
