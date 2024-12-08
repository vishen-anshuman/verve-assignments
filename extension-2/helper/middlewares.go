package helperfunc

import (
	"extension-2/app"
	"fmt"
	"net/http"
)

// "Before" aspect: This function checks if the ID is already being processed.
func BeforeProcessing(r *http.Request) int {
	query := r.URL.Query()
	idParam := query.Get("id")
	if idParam == "" {
		return http.StatusBadRequest
	}
	if err := idProcessingMiddleware(idParam); err != nil {
		return http.StatusConflict
	}
	return 0
}

// "After" aspect: This function cleans up the inProgress and UniqueIDCache after the handler logic.
func AfterProcessing(idParam string) {
	appInst := app.GetAppConst()
	idParamKey := fmt.Sprintf(app.PROCESSING_ID_FORMAT, idParam)
	appInst.Mu.Lock()
	appInst.RedisService.DeleteCache(idParamKey)
	appInst.Mu.Unlock()
}

func idProcessingMiddleware(idParam string) error {
	appInst := app.GetAppConst()
	appInst.Mu.Lock()
	idParamKey := fmt.Sprintf(app.PROCESSING_ID_FORMAT, idParam)
	if value, _ := appInst.RedisService.ReadFromCache(idParamKey); value != "" {
		appInst.Mu.Unlock()
		return fmt.Errorf("ID is already being processed")
	}
	appInst.RedisService.WriteToCache(idParamKey, "exists")
	appInst.Mu.Unlock()
	return nil
}
