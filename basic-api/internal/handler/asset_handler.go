package handler

import (
	"encoding/json"
	"errors"
	"mini-asm/internal/model"
	"mini-asm/internal/service"
	"net/http"
	"strings"
)

// AssetHandler handles HTTP requests for asset operations
// It's responsible for HTTP concerns only (parsing, status codes, JSON)
type AssetHandler struct {
	service *service.AssetService
}

// NewAssetHandler creates a new asset handler
func NewAssetHandler(service *service.AssetService) *AssetHandler {
	return &AssetHandler{
		service: service,
	}
}

// CreateAssetRequest represents the request body for creating an asset
type CreateAssetRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status,omitempty"` // Optional
}

// CreateAssetBatchResponse represents the response for a batch creation
type CreateAssetBatchResponse struct {
	Assets []model.Asset `json:"assets"`
	Count  int           `json:"count"`
}

// DeleteAssetBatchResponse represents the response for a batch deletion
type DeleteAssetBatchResponse struct {
	Deleted  int `json:"deleted"`
	NotFound int `json:"not_found"`
}

// UpdateAssetRequest represents the request body for updating an asset
type UpdateAssetRequest struct {
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	Status string `json:"status,omitempty"`
}

// CreateAsset handles POST /assets
func (h *AssetHandler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	var req CreateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	asset, err := h.service.CreateAsset(req.Name, req.Type)
	if err != nil {
		statusCode := mapErrorToStatus(err)
		RespondError(w, statusCode, err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, asset)
}

// CreateAssetBatchRequest represents the request body for batch creating assets
type CreateAssetBatchRequest struct {
	Assets []CreateAssetRequest `json:"assets"`
}

// BatchCreateAssets handles POST /assets/batch
func (h *AssetHandler) BatchCreateAssets(w http.ResponseWriter, r *http.Request) {
	var req CreateAssetBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Limit to 100 assets
	if len(req.Assets) > 100 {
		RespondError(w, http.StatusBadRequest, "Maximum 100 assets allowed per batch")
		return
	}

	batchItems := make([]struct{ Name, Type string }, len(req.Assets))
	for i, a := range req.Assets {
		batchItems[i] = struct{ Name, Type string }{Name: a.Name, Type: a.Type}
	}

	assets, err := h.service.CreateAssetsBatch(r.Context(), batchItems)
	if err != nil {
		statusCode := mapErrorToStatus(err)
		RespondError(w, statusCode, err.Error())
		return
	}

	// Extract IDs for the response
	ids := make([]string, len(assets))
	for i, asset := range assets {
		ids[i] = asset.ID
	}

	// Match expected homework response
	RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"created": len(assets),
		"ids":     ids,
	})
}

// ListAssets handles GET /assets
func (h *AssetHandler) ListAssets(w http.ResponseWriter, r *http.Request) {
	// Get query parameters for filtering/searching
	assetType := r.URL.Query().Get("type")
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")

	var assets []*model.Asset
	var err error

	// Determine which operation to perform
	if search != "" {
		// Search by name
		assets, err = h.service.SearchAssets(search)
	} else if assetType != "" || status != "" {
		// Filter by type and/or status
		assets, err = h.service.FilterAssets(assetType, status)
	} else {
		// Get all assets
		assets, err = h.service.GetAllAssets()
	}

	if err != nil {
		statusCode := mapErrorToStatus(err)
		RespondError(w, statusCode, err.Error())
		return
	}

	// Return empty array instead of null if no assets
	if assets == nil {
		assets = []*model.Asset{}
	}

	RespondJSON(w, http.StatusOK, assets)
}

// GetAsset handles GET /assets/{id}
func (h *AssetHandler) GetAsset(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL path
	id := r.PathValue("id")
	if id == "" {
		RespondError(w, http.StatusBadRequest, "Asset ID is required")
		return
	}

	// Call service
	asset, err := h.service.GetAssetByID(id)
	if err != nil {
		statusCode := mapErrorToStatus(err)
		RespondError(w, statusCode, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, asset)
}

// UpdateAsset handles PUT /assets/{id}
func (h *AssetHandler) UpdateAsset(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL path
	id := r.PathValue("id")
	if id == "" {
		RespondError(w, http.StatusBadRequest, "Asset ID is required")
		return
	}

	// Parse request body
	var req UpdateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Call service
	asset, err := h.service.UpdateAsset(id, req.Name, req.Type, req.Status)
	if err != nil {
		statusCode := mapErrorToStatus(err)
		RespondError(w, statusCode, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, asset)
}

// DeleteAsset handles DELETE /assets/{id}
func (h *AssetHandler) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL path
	id := r.PathValue("id")
	if id == "" {
		RespondError(w, http.StatusBadRequest, "Asset ID is required")
		return
	}

	// Call service
	if err := h.service.DeleteAsset(id); err != nil {
		statusCode := mapErrorToStatus(err)
		RespondError(w, statusCode, err.Error())
		return
	}

	// 204 No Content - successful deletion
	w.WriteHeader(http.StatusNoContent)
}

// BatchDeleteAssets handles DELETE /assets/batch?ids=id1,id2
func (h *AssetHandler) BatchDeleteAssets(w http.ResponseWriter, r *http.Request) {
	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		RespondError(w, http.StatusBadRequest, "ids parameter is required")
		return
	}

	ids := strings.Split(idsParam, ",")
	deleted, notFound, err := h.service.DeleteAssetsBatch(r.Context(), ids)
	if err != nil {
		statusCode := mapErrorToStatus(err)
		RespondError(w, statusCode, err.Error())
		return
	}

	// Return homework expected response
	RespondJSON(w, http.StatusOK, map[string]int{
		"deleted":   deleted,
		"not_found": notFound,
	})
}

// GetStats handles GET /assets/stats
// Returns total asset count and breakdown by type and status
func (h *AssetHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStats()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, stats)
}

// CountAssets handles GET /assets/count
// Supports optional query params: type, status
func (h *AssetHandler) CountAssets(w http.ResponseWriter, r *http.Request) {
	assetType := r.URL.Query().Get("type")
	status := r.URL.Query().Get("status")

	count, err := h.service.CountByFilter(assetType, status)
	if err != nil {
		statusCode := mapErrorToStatus(err)
		RespondError(w, statusCode, err.Error())
		return
	}

	// Build filters map for response (only include provided filters)
	filters := map[string]string{}
	if assetType != "" {
		filters["type"] = assetType
	}
	if status != "" {
		filters["status"] = status
	}

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"count":   count,
		"filters": filters,
	})
}

// mapErrorToStatus maps service layer errors to HTTP status codes
func mapErrorToStatus(err error) int {
	switch {
	case errors.Is(err, model.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, model.ErrInvalidInput),
		errors.Is(err, model.ErrEmptyName),
		errors.Is(err, model.ErrInvalidType),
		errors.Is(err, model.ErrInvalidStatus):
		return http.StatusBadRequest
	case errors.Is(err, model.ErrDuplicate):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// RespondJSON writes a JSON response
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log error but can't change status code (already written)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// RespondError writes a JSON error response
func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]string{
		"error": message,
	})
}

/*
🎓 NOTES:

1. Handler Responsibilities:
   ✅ Parse HTTP request (JSON, query params, URL params)
   ✅ Call appropriate service method
   ✅ Map errors to HTTP status codes
   ✅ Format HTTP response (JSON)
   ❌ Business logic
   ❌ Validation (except HTTP-specific)
   ❌ Data access

2. HTTP Request Parsing:

   Body (JSON):
   var req CreateAssetRequest
   json.NewDecoder(r.Body).Decode(&req)

   Query params:
   type := r.URL.Query().Get("type")

   URL path params:
   id := r.PathValue("id")  // Go 1.22+

3. HTTP Status Codes:
   200 OK           - Successful GET/PUT
   201 Created      - Successful POST
   204 No Content   - Successful DELETE
   400 Bad Request  - Invalid input
   404 Not Found    - Resource doesn't exist
   409 Conflict     - Duplicate
   500 Internal     - Server error

4. Error Mapping:
   func mapErrorToStatus(err error) int {
       switch {
       case errors.Is(err, model.ErrNotFound):
           return 404
       case errors.Is(err, model.ErrInvalidInput):
           return 400
       ...
       }
   }

   Q: Tại sao không return status code từ service?
   A: Service layer không biết về HTTP!
      Có thể reuse service cho CLI, gRPC, etc.

5. JSON Response Helpers:
   RespondJSON() - generic
   RespondError() - consistent error format

   {"error": "message"} - standard format

6. Request/Response Structs:
   type CreateAssetRequest struct {
       Name string `json:"name"`
       Type string `json:"type"`
   }

   Q: Tại sao không dùng model.Asset trực tiếp?
   A: API request != domain model
      - Request có thể có extra fields (passwords, etc.)
      - Response có thể exclude fields (sensitive data)
      - Clear API contract

7. Query Parameters:
   GET /assets?type=domain&status=active
   → r.URL.Query().Get("type")

   Flexible filtering!

8. Go 1.22+ Path Values:
   Pattern: "GET /assets/{id}"
   Get value: r.PathValue("id")

   Trước Go 1.22: phải dùng regex hoặc third-party router

📝 COMMON MISTAKES:

❌ Mistake 1: Business logic trong handler
func (h *Handler) CreateAsset(w, r) {
    // Parse request
    if req.Name == "" { return } // Validation ở đây - WRONG!
    asset.ID = uuid.New()          // Business logic - WRONG!
}
→ Should be in service layer!

❌ Mistake 2: SQL trong handler
func (h *Handler) CreateAsset(w, r) {
    db.Exec("INSERT INTO...")  // Data access - WRONG!
}
→ Should be in storage layer!

❌ Mistake 3: Not checking errors
json.NewDecoder(r.Body).Decode(&req)  // No error check - WRONG!

✅ Always check errors and respond appropriately

🔄 REQUEST FLOW EXAMPLE:

Client sends:
POST /assets
{"name": "example.com", "type": "domain"}

Handler:
1. Parse JSON → CreateAssetRequest
2. Call service.CreateAsset("example.com", "domain")
3. Service returns (*Asset, nil) or (nil, error)
4. Map result to HTTP response

Success response:
201 Created
{
  "id": "uuid",
  "name": "example.com",
  "type": "domain",
  "status": "active",
  "created_at": "2026-03-03T10:00:00Z",
  "updated_at": "2026-03-03T10:00:00Z"
}

Error response:
400 Bad Request
{
  "error": "name is required"
}

❓ QUESTIONS TO ASK:

1. Tại sao cần RespondError helper function?
   → Consistency, DRY principle

2. Handler có nên log không?
   → Có! (Buổi 5 sẽ add logging middleware)

3. Làm sao handle CORS?
   → Middleware! (Buổi 6)

4. PUT vs PATCH?
   → PUT = replace entire resource
   → PATCH = partial update
   → We use PUT with partial update logic
*/
