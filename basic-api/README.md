# 🔨 Buổi 2: Basic API Development

## Mục Tiêu

- ✅ Implement Clean Architecture với 4 layers
- ✅ CRUD operations: Create, Read (List, Get by ID)
- ✅ In-memory storage với thread-safe
- ✅ Hiểu dependency injection và interfaces

## So Sánh với Buổi 1

| Buổi 1          | Buổi 2              |
| --------------- | ------------------- |
| 1 file main.go  | 4 layers riêng biệt |
| Hello World     | Full CRUD API       |
| No data storage | In-memory storage   |
| Monolithic      | Clean Architecture  |

## Architecture Overview

```
Request Flow:
HTTP Request
  → Handler (parse JSON, HTTP concerns)
    → Service (business logic, validation)
      → Storage (data persistence)
        → Model (domain entities)
```

## Key Changes

### 1. Entity Layer (`internal/model/`)

**New Files:**

- `asset.go` - Asset struct và constants
- `errors.go` - Custom error types

**Key Points:**

- Pure domain logic, no dependencies
- Struct tags cho JSON marshalling
- Constants for type safety

### 2. Storage Layer (`internal/storage/`)

**New Files:**

- `storage.go` - Interface definition
- `memory/memory.go` - In-memory implementation

**Key Points:**

- Interface cho flexibility (swap implementations)
- Thread-safety với sync.RWMutex
- Why interface? → Buổi 3 sẽ swap sang database!

### 3. Service Layer (`internal/service/`)

**New Files:**

- `asset_service.go` - Business logic
- `service.go` - Service interface

**Teaching Points:**

- Validation logic
- UUID generation
- Business rules (default status = active)
- Dependency injection (nhận Storage interface)

### 4. Handler Layer (`internal/handler/`)

**New Files:**

- `asset_handler.go` - HTTP handlers cho assets
- `health_handler.go` - Health check (refactored từ main)
- `handler.go` - Handler registry

**Key Points:**

- HTTP-specific code only
- JSON parsing và encoding
- Status codes (201, 400, 404, 500)
- Dependency injection (nhận Service)

### 5. Main (`cmd/server/main.go`)

**Changes:**

- Wire up all dependencies
- Register routes
- Remove business logic (moved to layers)

## API Endpoints

| Method | Path         | Description      | Status      |
| ------ | ------------ | ---------------- | ----------- |
| GET    | /health      | Health check     | ✅          |
| POST   | /assets      | Create asset     | ✅          |
| GET    | /assets      | List all assets  | ✅          |
| GET    | /assets/{id} | Get single asset | ✅          |
| PUT    | /assets/{id} | Update asset     | 🔜 Homework |
| DELETE | /assets/{id} | Delete asset     | 🔜 Homework |

## Testing

### 1. Health Check

```bash
curl http://localhost:8080/health
```

### 2. Create Asset

```bash
curl -X POST http://localhost:8080/assets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "example.com",
    "type": "domain"
  }'
```

### 3. List Assets

```bash
curl http://localhost:8080/assets
```

### 4. Get Single Asset

```bash
# Replace <id> with actual UUID from create response
curl http://localhost:8080/assets/<id>
```

**Summary**

- Validation BEFORE business logic
- UUID auto-generation
- Default values (status = active)
- Timestamps auto-set
- Service doesn't know HOW data is stored (memory? DB?)

- Handler only handles HTTP concerns
- No business logic here!
- Status codes: 201 for created, 400 for bad request
- Helper functions: RespondJSON, RespondError

- Each layer has single responsibility
- Easy to test each layer independently
- Easy to swap storage implementation
- Clear separation of concerns

## Resources

- Review: CLEAN_ARCHITECTURE.MD sections 2-4
- [Go Interfaces](https://go.dev/tour/methods/9)
- [Dependency Injection in Go](https://blog.drewolson.org/dependency-injection-in-go)
- [UUID package](https://github.com/google/uuid)
