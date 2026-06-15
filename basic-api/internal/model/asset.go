package model

import "time"

// Asset represents a public-facing resource (domain, IP, or service)
// This is our core domain entity - no dependencies on other layers
type Asset struct {
	ID        string    `json:"id"`         // UUID
	Name      string    `json:"name"`       // e.g., "example.com", "192.168.1.1"
	Type      string    `json:"type"`       // domain, ip, or service
	Status    string    `json:"status"`     // active or inactive
	CreatedAt time.Time `json:"created_at"` // Auto-set on creation
	UpdatedAt time.Time `json:"updated_at"` // Auto-updated
}

// Asset types - using constants for type safety
const (
	TypeDomain  = "domain"
	TypeIP      = "ip"
	TypeService = "service"
)

// Asset statuses
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

// AssetStats represents statistics about assets
type AssetStats struct {
	Total    int            `json:"total"`
	ByType   map[string]int `json:"by_type"`
	ByStatus map[string]int `json:"by_status"`
}

// IsValidType checks if the given type is valid
func IsValidType(t string) bool {
	return t == TypeDomain || t == TypeIP || t == TypeService
}

// IsValidStatus checks if the given status is valid
func IsValidStatus(s string) bool {
	return s == StatusActive || s == StatusInactive
}

/*
🎓 NOTES:

1. Pure Domain Entity:
   - No database tags (gorm, sql, etc.)
   - No HTTP concerns
   - Just business concepts
   - This is the "Entity Layer" in Clean Architecture

2. Struct Tags:
   - `json:"id"` - định nghĩa tên field trong JSON response
   - Quan trọng: nếu không có tag, Go sẽ export field name as-is
   - Example: ID → "ID" vs id → "id"

3. Constants vs Strings:
   - ✅ TypeDomain - compiler checked, typo-safe
   - ❌ "domain" - runtime error if typo
   - Best practice: use constants!

4. Helper Functions:
   - IsValidType(), IsValidStatus()
   - Used by service layer for validation
   - Keep validation logic reusable

5. Why time.Time?
   - Built-in JSON marshalling
   - Timezone aware
   - Easy comparison and manipulation


*/
