package storage

import (
	"context"
	"mini-asm/internal/model"
)

// Storage defines the interface for data access operations
// This is the contract that any storage implementation must fulfill
//
// Why interface?
// - Allows multiple implementations (memory, postgres, mongodb, etc.)
// - Makes testing easy (can mock)
// - Follows Dependency Inversion Principle
type Storage interface {
	// Create adds a new asset to storage
	Create(asset *model.Asset) error

	// BatchCreate adds multiple assets to storage safely
	BatchCreate(ctx context.Context, assets []model.Asset) error

	// GetAll retrieves all assets
	// Returns slice of pointers for efficiency (no copying)
	GetAll() ([]*model.Asset, error)

	// GetByID retrieves a single asset by its ID
	// Returns ErrNotFound if asset doesn't exist
	GetByID(id string) (*model.Asset, error)

	// Update modifies an existing asset
	// Returns ErrNotFound if asset doesn't exist
	Update(id string, asset *model.Asset) error

	// Delete removes an asset from storage
	// Returns ErrNotFound if asset doesn't exist
	Delete(id string) error

	// BatchDelete removes multiple assets from storage
	BatchDelete(ctx context.Context, ids []string) (deleted int, notFound int, err error)

	// Filter returns assets matching the given criteria
	// Empty string parameters are ignored (match all)
	Filter(assetType, status string) ([]*model.Asset, error)

	// Search finds assets by partial name match
	Search(query string) ([]*model.Asset, error)

	// GetStats returns aggregate statistics about all assets
	GetStats() (*model.AssetStats, error)

	// CountByFilter returns count of assets matching given type and/or status
	CountByFilter(assetType, status string) (int, error)

	// GetCount returns the total number of assets in storage
	GetCount() (int, error)
}

/*
🎓 NOTES:

1. Interface Design:
   - Define behavior, not implementation
   - Methods should be atomic and clear
   - Return errors, don't panic

2. Why Pointers?
   - []*model.Asset vs []model.Asset
   - Pointers avoid copying large structs
   - Allows modification through reference
   - More memory efficient for large datasets

3. Error Handling:
   - Return error as last value
   - Use model.ErrNotFound for consistency
   - Caller decides how to handle

4. Method Signatures:
   - Create(asset *model.Asset) - pointer: will be modified (ID, timestamps)
   - GetByID(id string) - string: immutable lookup
   - Filter/Search - flexible parameters

5. Interface Benefits:

   Buổi 2: MemoryStorage implements this
   type MemoryStorage struct { ... }
   func (m *MemoryStorage) Create(asset *model.Asset) error { ... }

   Buổi 3: PostgresStorage implements the SAME interface
   type PostgresStorage struct { ... }
   func (p *PostgresStorage) Create(asset *model.Asset) error { ... }

   Service layer doesn't change!
   type AssetService struct {
       storage Storage  // Works with ANY implementation!
   }

6. Testing Benefits:
   type MockStorage struct { ... }
   func (m *MockStorage) Create(asset *model.Asset) error {
       return nil // or test-specific behavior
   }

📝 COMPARISON:

Without Interface (BAD):
    type AssetService struct {
        storage *MemoryStorage  // Coupled to specific implementation
    }
    // Can't swap to database without changing service!

With Interface (GOOD):
    type AssetService struct {
        storage Storage  // Any implementation works
    }
    // Easy to swap implementations!
*/
