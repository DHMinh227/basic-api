package handler

import (
	"net/http"
	"time"

	"mini-asm/internal/service"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	startTime    time.Time
	assetService *service.AssetService
}

// NewHealthHandler creates a new health check handler
// Accepts an AssetService to report in-memory storage stats
func NewHealthHandler(assetService *service.AssetService) *HealthHandler {
	return &HealthHandler{
		startTime:    time.Now(),
		assetService: assetService,
	}
}

// StorageInfo holds the storage-related health info
type StorageInfo struct {
	Type       string `json:"type"`
	AssetCount int    `json:"asset_count"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status        string      `json:"status"`
	Storage       StorageInfo `json:"storage"`
	UptimeSeconds float64     `json:"uptime_seconds"`
	Timestamp     time.Time   `json:"timestamp"`
}

// Check handles GET /health
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	count, _ := h.assetService.GetCount()

	response := HealthResponse{
		Status: "ok",
		Storage: StorageInfo{
			Type:       "in-memory",
			AssetCount: count,
		},
		UptimeSeconds: time.Since(h.startTime).Seconds(),
		Timestamp:     time.Now(),
	}

	RespondJSON(w, http.StatusOK, response)
}

/*
🎓 NOTES:

Refactored from Session 1:
- Session 1: Health check logic in main.go
- Session 2: Extracted to separate handler
- Homework Task 5: Added storage info (asset_count, type, uptime_seconds)

Benefits:
- Consistent with other handlers
- Can add more health checks (database, etc.) in Session 3
- Reusable and testable
*/
