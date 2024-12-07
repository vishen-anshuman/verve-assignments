package helperfunc

import (
	"extension-2/redisservice"
	"fmt"
	"net/http"
	"sync"
)

// "Before" aspect: This function checks if the ID is already being processed.
func BeforeProcessing(r *http.Request) int {
	query := r.URL.Query()
	idParam := query.Get("id")
	if idParam == "" {
		return http.StatusConflict
	}
	if err := idProcessingMiddleware(idParam); err != nil {
		return http.StatusBadRequest
	}
	return 0
}

// "After" aspect: This function cleans up the inProgress and UniqueIDCache after the handler logic.
func AfterProcessing(idParam string) {
	var mu sync.Mutex
	mu.Lock()
	redisservice.DeleteCache(idParam)
	mu.Unlock()
}

func idProcessingMiddleware(idParam string) error {
	var mu sync.Mutex
	mu.Lock()
	if value, _ := redisservice.ReadFromCache(idParam); value != "" {
		mu.Unlock()
		return fmt.Errorf("ID is already being processed")
	}
	redisservice.WriteToCache(idParam, "exists")
	mu.Unlock()
	return nil
}
