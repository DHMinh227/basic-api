package memory

import (
	"context"
	"mini-asm/internal/model"
	"sort"
	"strings"
	"sync"
)

// MemoryStorage implements the Storage interface using in-memory map
// This is a simple implementation for development and testing
// Data is lost when the application restarts
type MemoryStorage struct {
	mu   sync.RWMutex            // protects concurrent access
	data map[string]*model.Asset // key = asset ID, value = asset pointer
}

// NewMemoryStorage creates a new in-memory storage instance
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]*model.Asset),
	}
}

// Create adds a new asset to memory
// Thread-safe with write lock
func (m *MemoryStorage) Create(asset *model.Asset) error {
	// ID assignment handled in service layer before lock
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[asset.ID]; exists {
		return model.ErrDuplicate
	}

	// Store a copy to prevent external modifications from affecting storage
	assetCopy := *asset
	m.data[asset.ID] = &assetCopy
	return nil
}

// BatchCreate adds multiple assets to storage safely
func (m *MemoryStorage) BatchCreate(ctx context.Context, assets []model.Asset) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// First pass: check for duplicates
	for _, asset := range assets {
		if _, exists := m.data[asset.ID]; exists {
			return model.ErrDuplicate
		}
	}

	// Second pass: store copies
	for _, asset := range assets {
		assetCopy := asset
		m.data[asset.ID] = &assetCopy
	}

	return nil
}

// GetAll returns all assets sorted by creation time (newest first)
// Thread-safe with read lock
func (m *MemoryStorage) GetAll() ([]*model.Asset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	assets := make([]*model.Asset, 0, len(m.data))
	for _, asset := range m.data {
		// Return copies to prevent callers from modifying storage state directly
		assetCopy := *asset
		assets = append(assets, &assetCopy)
	}

	// Sort by created_at descending (newest first)
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].CreatedAt.After(assets[j].CreatedAt)
	})

	return assets, nil
}

// GetByID retrieves a single asset by ID
// Thread-safe with read lock
func (m *MemoryStorage) GetByID(id string) (*model.Asset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	asset, exists := m.data[id]
	if !exists {
		return nil, model.ErrNotFound
	}

	// Return a copy
	assetCopy := *asset
	return &assetCopy, nil
}

// Update modifies an existing asset
// Thread-safe with write lock
func (m *MemoryStorage) Update(id string, asset *model.Asset) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[id]; !exists {
		return model.ErrNotFound
	}

	// Store a copy
	assetCopy := *asset
	m.data[id] = &assetCopy
	return nil
}

// Delete removes an asset from storage
// Thread-safe with write lock
func (m *MemoryStorage) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[id]; !exists {
		return model.ErrNotFound
	}

	delete(m.data, id)
	return nil
}

// BatchDelete removes multiple assets from storage
func (m *MemoryStorage) BatchDelete(ctx context.Context, ids []string) (deleted int, notFound int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, id := range ids {
		if _, exists := m.data[id]; exists {
			delete(m.data, id)
			deleted++
		} else {
			notFound++
		}
	}

	return deleted, notFound, nil
}

// Filter returns assets matching the given type and/or status
// Empty string parameters are ignored (match all)
// Thread-safe with read lock
func (m *MemoryStorage) Filter(assetType, status string) ([]*model.Asset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	assets := make([]*model.Asset, 0)
	for _, asset := range m.data {
		// Check type filter
		if assetType != "" && asset.Type != assetType {
			continue
		}

		// Check status filter
		if status != "" && asset.Status != status {
			continue
		}

		// Store copy
		assetCopy := *asset
		assets = append(assets, &assetCopy)
	}

	// Sort by created_at descending
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].CreatedAt.After(assets[j].CreatedAt)
	})

	return assets, nil
}

// Search finds assets where name contains the query string (case-insensitive)
// Thread-safe with read lock
func (m *MemoryStorage) Search(query string) ([]*model.Asset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	query = strings.ToLower(query)
	assets := make([]*model.Asset, 0)

	for _, asset := range m.data {
		if strings.Contains(strings.ToLower(asset.Name), query) {
			// Store copy
			assetCopy := *asset
			assets = append(assets, &assetCopy)
		}
	}

	// Sort by created_at descending
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].CreatedAt.After(assets[j].CreatedAt)
	})

	return assets, nil
}

// GetStats returns aggregate statistics about all assets in memory
// Thread-safe with read lock
func (m *MemoryStorage) GetStats() (*model.AssetStats, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := &model.AssetStats{
		Total:    len(m.data),
		ByType:   make(map[string]int),
		ByStatus: make(map[string]int),
	}

	for _, asset := range m.data {
		stats.ByType[asset.Type]++
		stats.ByStatus[asset.Status]++
	}

	return stats, nil
}

// CountByFilter returns the number of assets matching the given type and/or status
// Empty string parameters are ignored (match all)
// Thread-safe with read lock
func (m *MemoryStorage) CountByFilter(assetType, status string) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, asset := range m.data {
		if assetType != "" && asset.Type != assetType {
			continue
		}
		if status != "" && asset.Status != status {
			continue
		}
		count++
	}

	return count, nil
}

// GetCount returns the total number of assets currently in storage
// Thread-safe with read lock
func (m *MemoryStorage) GetCount() (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.data), nil
}

/*
🎓 NOTES:

1. Thread Safety với sync.RWMutex:

   Q: Tại sao cần mutex?
   A: Multiple HTTP requests có thể access storage cùng lúc
      → Data race → undefined behavior

   Q: RWMutex vs Mutex?
   A: RWMutex cho phép:
      - Multiple readers cùng lúc (RLock)
      - Single writer (Lock) - blocks all reads/writes
      Performance tốt hơn cho read-heavy workloads

   Demo (if time permits):
   go run -race cmd/server/main.go
   → Chạy concurrent requests
   → Without mutex: race detector warnings
   → With mutex: safe!

2. Lock Pattern:
   m.mu.Lock()
   defer m.mu.Unlock()  // Always unlock, even if panic

   Q: Tại sao defer?
   A: Đảm bảo unlock ngay cả khi có error/panic
      Prevent deadlock

3. Map Iteration:
   - for _, asset := range m.data
   - Order không guaranteed!
   - Phải sort nếu cần order

4. Sorting:
   sort.Slice(assets, func(i, j int) bool {
       return assets[i].CreatedAt.After(assets[j].CreatedAt)
   })
   - After() → descending (newest first)
   - Before() → ascending (oldest first)

5. Performance Considerations:
   - map[string]*Asset → O(1) lookup
   - GetAll() → O(n) + O(n log n) sort
   - Filter/Search → O(n) iteration

   Buổi 3 với database: có thể optimize với indexes!

6. Memory Management:
   - Pointers: không copy struct
   - make(map): allocate memory
   - Data lost on restart → need persistence!

📝 BUỔI 2 vs BUỔI 3:

Buổi 2: MemoryStorage
+ Pros: Fast, simple, no dependencies
- Cons: Data lost on restart, single instance only

Buổi 3: PostgresStorage
+ Pros: Persistent, scalable, multiple instances
- Cons: Slower, need database setup

Code changes needed: CHỈ 1 DÒNG trong main.go!
store := memory.NewMemoryStorage()  → store := postgres.NewPostgresStorage(db)

Service/Handler/Model: KHÔNG THAY ĐỔI!
→ This is the power of Clean Architecture!

🔍 CODE WALKTHROUGH TIPS:

1. Start with NewMemoryStorage - simple constructor
2. Show Create - write lock, duplicate check
3. Show GetAll - read lock, sorting
4. Explain RLock vs Lock difference
5. Demo concurrent access safety
6. Compare with Buổi 3 postgres implementation

❓ QUESTIONS TO ASK STUDENTS:

1. Tại sao return []*Asset chứ không phải []Asset?
   → Pointers, efficiency

2. Nếu không có mutex thì sao?
   → Data race, crash, corrupt data

3. Filter vs Search khác nhau như thế nào?
   → Filter: exact match, Search: partial match

4. Làm sao thêm Pagination?
   → Slice with offset/limit (homework!)
*/
