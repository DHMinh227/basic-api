package model

import "errors"

// Custom error types for domain logic
// These errors can be mapped to appropriate HTTP status codes by handlers
var (
	// ErrNotFound indicates the requested resource doesn't exist
	ErrNotFound = errors.New("asset not found")

	// ErrInvalidInput indicates validation failed
	ErrInvalidInput = errors.New("invalid input")

	// ErrDuplicate indicates trying to create a duplicate resource
	ErrDuplicate = errors.New("asset already exists")

	// ErrEmptyName indicates name field is required
	ErrEmptyName = errors.New("name is required")

	// ErrInvalidType indicates invalid asset type
	ErrInvalidType = errors.New("invalid asset type: must be domain, ip, or service")

	// ErrInvalidStatus indicates invalid asset status
	ErrInvalidStatus = errors.New("invalid status: must be active or inactive")
)

/*
🎓 NOTES:

1. Custom Error Types:
   - Define domain-specific errors
   - Better than generic error strings
   - Can be checked with errors.Is()

2. Error Handling Strategy:
   - Model layer: define error types
   - Service layer: return these errors
   - Handler layer: map to HTTP status codes
     * ErrNotFound → 404
     * ErrInvalidInput → 400
     * Others → 500

3. Why Global Variables?
   - Sentinel errors (can compare with ==)
   - Reusable across layers
   - Type-safe

Example usage:

Service layer:
    if name == "" {
        return ErrEmptyName
    }

Handler layer:
    if errors.Is(err, model.ErrNotFound) {
        return 404
    }

4. Alternative Approaches:
   - Custom error structs with more context
   - Wrapped errors with errors.Wrap()
   - Error codes (int constants)

   For this training, simple errors.New() is enough!

📝 BEST PRACTICES:

✅ Define errors in model package (domain layer)
✅ Descriptive error messages
✅ Consistent naming: Err* prefix

❌ Don't create error inline: errors.New("not found") everywhere
❌ Don't include HTTP details in error message
❌ Don't panic for expected errors
*/
