package helperfunc

import (
	"extension-2/app"
	"extension-2/redisservice"
	"fmt"
	"net/http"
)

// "Before" aspect: This function checks if the ID is already being processed.
func BeforeProcessing(r *http.Request, appInst *app.App) int {
	query := r.URL.Query()
	idParam := query.Get("id")
	if idParam == "" {
		return http.StatusBadRequest
	}
	if err := idProcessingMiddleware(idParam, appInst); err != nil {
		return http.StatusConflict
	}
	return 0
}

// "After" aspect: This function cleans up the inProgress and UniqueIDCache after the handler logic.
func AfterProcessing(idParam string, appInst *app.App) {
	appInst.Mu.Lock()
	redisservice.DeleteCache(idParam)
	appInst.Mu.Unlock()
}

func idProcessingMiddleware(idParam string, appInst *app.App) error {
	appInst.Mu.Lock()
	if value, _ := redisservice.ReadFromCache(idParam); value != "" {
		appInst.Mu.Unlock()
		return fmt.Errorf("ID is already being processed")
	}
	redisservice.WriteToCache(idParam, "exists")
	appInst.Mu.Unlock()
	return nil
}
