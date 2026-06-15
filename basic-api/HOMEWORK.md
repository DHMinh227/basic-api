# 📝 Bài Tập Về Nhà - Sessions 1-2

**Deadline:** Trước ngày thứ 3  
**Cách nộp:** Push lên Git repository cá nhân tại branch homework. Mời dinhmanhtan (dmtangtnd@gmail.com) vào project. Tạo pull request từ branch homework vào main. Set reviewer là dinhmanhtan  
**Note**: Project có thể viết bằng ngôn ngữ khác, không bắt buộc phải dùng Go. Nếu sử dụng ngôn ngữ khác cần mô tả cách cài đặt và chạy project.

---

## 📑 Mục Lục

- [Yêu Cầu Chung](#yêu-cầu-chung)
- [📤 Nộp Bài](#-nộp-bài)
- [Bài 1: Statistics APIs (20 điểm)](#bài-1-statistics-apis-20-điểm)
  - [1.1 Get Assets Statistics](#11-get-assets-statistics)
  - [1.2 Count Assets by Filter](#12-count-assets-by-filter)
- [Bài 2: Batch Create Assets (25 điểm)](#bài-2-batch-create-assets-25-điểm)
- [Bài 3: Batch Delete Assets (20 điểm)](#bài-3-batch-delete-assets-20-điểm)
- [Bài 4: Concurrent-safe Create (25 điểm) ⭐](#bài-4-concurrent-safe-create-25-điểm-)
- [Bài 5: In-memory Health Check (15 điểm)](#bài-5-in-memory-health-check-15-điểm)
- [Bài 6: Pagination & Filtering (15 điểm) - BONUS 🌟](#bài-6-pagination--filtering-15-điểm---bonus-)
- [Bài 7: Search by Name (10 điểm) - BONUS 🌟](#bài-7-search-by-name-10-điểm---bonus-)
- [📊 Chấm Điểm](#-chấm-điểm)
- [💡 Gợi Ý & Tips](#-gợi-ý--tips)
- [📚 Tài Liệu Tham Khảo](#-tài-liệu-tham-khảo)
- [🚀 Bonus Challenges](#-bonus-challenges)

---

## Yêu Cầu Chung

- ✅ Code phải chạy được, không có lỗi.
- ✅ Follow Clean Architecture như đã học (handler → service → storage → model).
- ✅ Dùng **in-memory storage** (không dùng database cho Day 1).
- ✅ Có error handling đầy đủ.
- ✅ Test được bằng curl hoặc Postman.

---

## 📤 Nộp Bài

### Cần nộp:

**File `SUBMISSION.md`** hoặc **`SUBMISSION.pdf`** gồm:

```markdown
# Homework Submission

**Họ tên:** [Tên của bạn]

## Các bài đã hoàn thành

- [x] Bài 1: Statistics APIs
- [x] Bài 2: Batch Create
- [x] Bài 3: Batch Delete
- [x] Bài 4: Concurrent-safe Create
- [x] Bài 5: In-memory Health Check
- [ ] Bài 6: Pagination (Bonus)
- [ ] Bài 7: Search (Bonus)
```

Mỗi bài cần 1 file test screenshots hoặc command outputs chứng minh.  
File này đặt trong thư mục [homeworks/submissions](../submissions/).

---

## Bài 1: Statistics APIs (20 điểm)

**Yêu cầu:** Implement API để lấy thống kê về assets từ in-memory storage.

### 1.1 Get Assets Statistics

- **Endpoint:** `GET /assets/stats`
- **Response:** 200 OK
  ```json
  {
    "total": 150,
    "by_type": {
      "domain": 100,
      "ip": 40,
      "service": 10
    },
    "by_status": {
      "active": 120,
      "inactive": 30
    }
  }
  ```

### 1.2 Count Assets by Filter

- **Endpoint:** `GET /assets/count`
- **Query params:** `type`, `status` (optional)
- **Response:** 200 OK
  ```json
  {
    "count": 85,
    "filters": {
      "type": "domain",
      "status": "active"
    }
  }
  ```

**Test:**

```bash
# Get statistics
curl http://localhost:8080/assets/stats

# Count all
curl http://localhost:8080/assets/count

# Count by type
curl "http://localhost:8080/assets/count?type=domain"

# Count by type and status
curl "http://localhost:8080/assets/count?type=domain&status=active"
```

---

## Bài 2: Batch Create Assets (25 điểm)

**Yêu cầu:** Tạo nhiều assets cùng lúc trong 1 request.

### API Specification

- **Endpoint:** `POST /assets/batch`
- **Request body:**
  ```json
  {
    "assets": [
      { "name": "domain1.com", "type": "domain" },
      { "name": "domain2.com", "type": "domain" },
      { "name": "192.168.1.1", "type": "ip" }
    ]
  }
  ```
- **Response:** 201 Created
  ```json
  {
    "created": 3,
    "ids": ["uuid-1", "uuid-2", "uuid-3"]
  }
  ```

### Yêu cầu kỹ thuật:

- Mô phỏng nguyên tắc **all or nothing** trong memory:
  - Validate toàn bộ list trước.
  - Chỉ insert khi tất cả items đều hợp lệ.
- Nếu 1 asset validation fail → không insert asset nào.
- Limit tối đa 100 assets/request.
- Validate từng asset trước khi insert.

**Test:**

```bash
# Success case
curl -X POST http://localhost:8080/assets/batch \
  -H "Content-Type: application/json" \
  -d '{
    "assets": [
      {"name":"test1.com","type":"domain"},
      {"name":"test2.com","type":"domain"}
    ]
  }'

# Error case (invalid type) - should create none
curl -X POST http://localhost:8080/assets/batch \
  -H "Content-Type: application/json" \
  -d '{
    "assets": [
      {"name":"test1.com","type":"domain"},
      {"name":"test2.com","type":"invalid_type"}
    ]
  }'
# Expected: 400 Bad Request, none created
```

---

## Bài 3: Batch Delete Assets (20 điểm)

**Yêu cầu:** Xóa nhiều assets cùng lúc.

### API Specification

- **Endpoint:** `DELETE /assets/batch`
- **Query params:** `?ids=uuid1,uuid2,uuid3`
- **Response:** 200 OK
  ```json
  {
    "deleted": 3,
    "not_found": 0
  }
  ```

### Behavior:

- Xóa tất cả IDs hợp lệ.
- Bỏ qua IDs không tồn tại (không trả lỗi).
- Return số lượng đã xóa và không tìm thấy.

**Test:**

```bash
# Create test assets first
ID1=$(curl -s -X POST http://localhost:8080/assets \
  -H "Content-Type: application/json" \
  -d '{"name":"test1.com","type":"domain"}' | jq -r '.id')

ID2=$(curl -s -X POST http://localhost:8080/assets \
  -H "Content-Type: application/json" \
  -d '{"name":"test2.com","type":"domain"}' | jq -r '.id')

# Batch delete (include 1 fake ID)
curl -X DELETE "http://localhost:8080/assets/batch?ids=$ID1,$ID2,fake-uuid-123"

# Expected response:
# {"deleted": 2, "not_found": 1}

# Verify deletion
curl http://localhost:8080/assets/$ID1
# Expected: 404 Not Found
```

---

## Bài 4: Concurrent-safe Create (25 điểm) ⭐

**Yêu cầu:** Đảm bảo in-memory storage an toàn khi nhiều request tạo asset cùng lúc.

### Specification:

- Dùng `sync.RWMutex` (hoặc cơ chế lock tương đương).
- Không bị race condition khi gọi concurrent create.
- Không tạo trùng ID.
- Server vẫn phản hồi ổn định khi có burst request.

### Gợi ý test nhanh:

```bash
# Bắn 20 request create song song
for i in $(seq 1 20); do
  curl -s -X POST http://localhost:8080/assets \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"concurrent-$i.com\",\"type\":\"domain\"}" &
done
wait

# Kiểm tra tổng số lượng
curl http://localhost:8080/assets/count
```

### Tiêu chí đạt:

- Không crash server.
- Không data corruption.
- Kết quả count phù hợp với số request hợp lệ.

---

## Bài 5: In-memory Health Check (15 điểm)

**Yêu cầu:** Nâng cấp `/health` endpoint để phản ánh trạng thái ứng dụng với storage in-memory.

### API Specification

- **Endpoint:** `GET /health`
- **Response:** 200 OK

  ```json
  {
    "status": "ok",
    "storage": {
      "type": "in-memory",
      "asset_count": 42
    },
    "uptime_seconds": 3600,
    "timestamp": "2026-03-06T10:00:00Z"
  }
  ```

### Implementation hints:

- `HealthHandler` có thể nhận thêm service/storage để lấy `asset_count`.
- Track thời điểm app start để tính `uptime_seconds`.
- Trả 200 khi server hoạt động bình thường.

**Test:**

```bash
# Health check
curl http://localhost:8080/health | jq

# Create assets rồi check lại
curl -X POST http://localhost:8080/assets \
  -H "Content-Type: application/json" \
  -d '{"name":"health-test.com","type":"domain"}'

curl http://localhost:8080/health | jq
# Expected: storage.asset_count tăng lên
```

---

## Bài 6: Pagination & Filtering (15 điểm) - BONUS 🌟

**Yêu cầu:** Thêm phân trang và filter cho list assets (thực hiện trên in-memory list).

### API Specification

- **Endpoint:** `GET /assets`
- **Query params:**
  - `page` (default: 1)
  - `limit` (default: 20, max: 100)
  - `type` (optional: domain, ip, service)
  - `status` (optional: active, inactive)

- **Response:**
  ```json
  {
    "data": [...],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "total_pages": 8
    }
  }
  ```

**Test:**

```bash
# Page 1, 10 items
curl "http://localhost:8080/assets?page=1&limit=10"

# Filter by type
curl "http://localhost:8080/assets?type=domain"

# Combine filters
curl "http://localhost:8080/assets?page=2&limit=20&type=domain&status=active"
```

---

## Bài 7: Search by Name (10 điểm) - BONUS 🌟

**Yêu cầu:** Tìm kiếm assets theo tên (partial match) trên in-memory data.

### API Specification

- **Endpoint:** `GET /assets/search`
- **Query params:** `q` (search query, required)
- **Response:** Array of matching assets (max 100)
- **Behavior:** Case-insensitive, partial match

**Test:**

```bash
# Search for "example"
curl "http://localhost:8080/assets/search?q=example"

# Search for ".com"
curl "http://localhost:8080/assets/search?q=.com"

# Case insensitive
curl "http://localhost:8080/assets/search?q=DOMAIN"
```

---

## 📊 Chấm Điểm

| Bài Tập                 | Điểm    | Bắt Buộc    |
| ----------------------- | ------- | ----------- |
| Bài 1: Statistics       | 20      | ✅ Bắt buộc |
| Bài 2: Batch Create     | 25      | ✅ Bắt buộc |
| Bài 3: Batch Delete     | 20      | ✅ Bắt buộc |
| Bài 4: Concurrent-safe  | 25      | ✅ Bắt buộc |
| Bài 5: In-memory Health | 15      | ✅ Bắt buộc |
| Bài 6: Pagination       | 15      | 🌟 Bonus    |
| Bài 7: Search           | 10      | 🌟 Bonus    |
| **Tổng bắt buộc**       | **105** |             |
| **Tổng có bonus**       | **130** |             |

---

## 💡 Gợi Ý & Tips

### Validate all before write (all or nothing trong memory):

```go
for _, in := range req.Assets {
    if err := validateAsset(in); err != nil {
        return nil, err // Stop here, do not write anything yet
    }
}

created := make([]Asset, 0, len(req.Assets))
for _, in := range req.Assets {
    created = append(created, toAsset(in))
}

if err := storage.BatchCreate(ctx, created); err != nil {
    return nil, err
}
```

### Filter trong memory:

```go
func match(a Asset, t, s string) bool {
    if t != "" && a.Type != t {
        return false
    }
    if s != "" && a.Status != s {
        return false
    }
    return true
}
```

### Parse query string IDs:

```go
idsParam := r.URL.Query().Get("ids")
if idsParam == "" {
    return nil, errors.New("ids parameter required")
}

ids := strings.Split(idsParam, ",")
// ids = ["uuid1", "uuid2", "uuid3"]
```

### Concurrency check bằng race detector:

```bash
go test ./... -race
```

---

## 📚 Tài Liệu Tham Khảo

- [Go map and concurrency](https://go.dev/blog/maps)
- [sync package](https://pkg.go.dev/sync)
- [Go race detector](https://go.dev/doc/articles/race_detector)
- [RESTful API Design](https://restfulapi.net/)

## 🚀 Bonus Challenges

1. **Rate Limiting:** Giới hạn số request/phút từ mỗi IP.
2. **Caching:** Cache kết quả list assets trong memory (5 phút).
3. **Audit Log:** Log mọi CREATE/UPDATE/DELETE ra file.
4. **Soft Delete:** Thêm `deleted_at` timestamp thay vì xóa hẳn.
5. **Import CSV:** Upload file CSV để tạo nhiều assets.
6. **Export CSV:** Download assets dưới dạng CSV.
7. **Webhooks:** Gọi webhook khi có asset mới được tạo.

---

**Chúc các bạn làm bài tốt! Có thắc mắc hỏi trên group nhé! 🚀**

---
