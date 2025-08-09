# Gau Kanban Service

## Mô tả dự án

Gau Kanban Service là một RESTful API service được xây dựng bằng Go để quản lý bảng Kanban. Dự án cung cấp các tính năng hoàn chỉnh để tạo và quản lý boards, columns, tickets, assignments và labels cho việc quản lý dự án theo phương pháp Kanban.

## Tính năng chính

### 🏗️ Quản lý Columns
- Tạo, sửa, xóa columns
- Sắp xếp lại vị trí columns
- Quản lý thứ tự hiển thị

### 🎫 Quản lý Tickets
- Tạo tickets với ticket number tự động (TASK-XXXX)
- CRUD operations cho tickets
- Di chuyển tickets giữa các columns
- Drag & drop với position management thông minh
- Tự động sắp xếp vị trí khi tạo ticket mới (luôn ở cuối column)
- Hỗ trợ due date và priority

### 👥 Quản lý Assignments
- Gán người dùng vào tickets
- Quản lý thông tin assignees (user_id, user_full_name)
- Xóa assignments theo user hoặc ticket
- Hiển thị assignees trong thông tin tickets

### 🏷️ Quản lý Labels
- Tạo và quản lý labels với màu sắc
- Gán labels vào tickets
- Quản lý many-to-many relationship

### 💬 Quản lý Comments
- Thêm comments vào tickets
- Quản lý discussions cho từng ticket

## Công nghệ sử dụng

- **Backend**: Go (Golang) với Gin framework
- **Database**: PostgreSQL với GORM ORM
- **Migration**: golang-migrate
- **Container**: Docker & Docker Compose
- **Architecture**: Clean Architecture với Repository pattern

## Cấu trúc dự án

```
gau-kanban-service/
├── main.go                 # Entry point
├── Dockerfile             # Docker configuration
├── entrypoint.sh          # Docker entrypoint script
├── config/                # Configuration management
├── controller/            # HTTP handlers
├── entity/               # Database models
├── repository/           # Data access layer
├── routes/               # API routes definition
├── migrations/           # Database migrations
├── infra/                # Infrastructure setup
└── utils/                # Utility functions
```

## API Endpoints

### Column Management
```
POST   /api/v2/kanban/columns              # Tạo column mới
GET    /api/v2/kanban/columns              # Lấy danh sách columns
PUT    /api/v2/kanban/columns/:id          # Cập nhật column
DELETE /api/v2/kanban/columns/:id          # Xóa column
PATCH  /api/v2/kanban/columns/:id/position # Cập nhật vị trí column
```

### Ticket Management
```
POST   /api/v2/kanban/tickets                    # Tạo ticket mới
GET    /api/v2/kanban/tickets                    # Lấy danh sách tickets
GET    /api/v2/kanban/tickets/:id                # Lấy thông tin ticket
PUT    /api/v2/kanban/tickets/:id                # Cập nhật ticket
DELETE /api/v2/kanban/tickets/:id                # Xóa ticket
PATCH  /api/v2/kanban/tickets/move               # Di chuyển ticket
PATCH  /api/v2/kanban/tickets/move-with-position # Di chuyển ticket với vị trí cụ thể
PATCH  /api/v2/kanban/tickets/:id/position       # Cập nhật vị trí ticket
```

### Assignment Management
```
POST   /api/v2/kanban/assignments                     # Tạo assignment
PUT    /api/v2/kanban/assignments/:id                 # Cập nhật assignment
DELETE /api/v2/kanban/assignments/:id                 # Xóa assignment
DELETE /api/v2/kanban/users/:user_id/assignments      # Xóa tất cả assignments của user
GET    /api/v2/kanban/tickets/:ticket_id/assignments  # Lấy assignments của ticket
```

### Kanban Board
```
GET    /api/v2/kanban/board       # Lấy toàn bộ kanban board
GET    /api/v2/kanban/tag-colors  # Lấy màu sắc tags
```

## Cài đặt và chạy

### Yêu cầu hệ thống
- Go 1.23+
- PostgreSQL 12+
- Docker & Docker Compose (optional)

### Chạy với Docker
```bash
# Clone repository
git clone <repository-url>
cd gau-kanban-service

# Chạy với Docker Compose
docker-compose up -d

# Service sẽ chạy trên port 8080
```

### Chạy development
```bash
# Cài đặt dependencies
go mod tidy

# Setup database (PostgreSQL)
# Tạo database: gau_kanban

# Chạy migrations
migrate -path migrations -database "postgres://username:password@localhost/gau_kanban?sslmode=disable" up

# Chạy service
go run main.go
```

### Environment Variables
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=gau_kanban
DB_SSLMODE=disable
PORT=8080
```

## Tính năng nổi bật

### 🎯 Smart Position Management
- Tự động sắp xếp vị trí tickets khi drag & drop
- Hỗ trợ di chuyển giữa các columns với transaction safety
- Tickets mới luôn được đặt ở cuối column

### 🔢 Auto Ticket Numbering
- Tự động tạo ticket number theo format TASK-XXXX
- Unique và sequential numbering

### 📊 Rich Data Response
- API responses bao gồm đầy đủ thông tin assignees
- Nested data cho kanban board view
- Optimized queries cho performance

### 🛡️ Data Integrity
- Database constraints và foreign keys
- Transaction handling cho complex operations
- Error handling và validation

## Database Schema

### Core Tables
- `columns`: Quản lý các cột kanban
- `tickets`: Quản lý các tickets/tasks
- `task_assignments`: Gán người dùng vào tickets
- `labels`: Quản lý labels/tags
- `ticket_labels`: Many-to-many relationship
- `ticket_comments`: Comments cho tickets

### Key Features
- UUID primary keys
- Timestamps tracking
- Position-based ordering
- Cascading deletes

## Migration Management

```bash
# Tạo migration mới
migrate create -ext sql -dir migrations -seq migration_name

# Chạy migrations
migrate -path migrations -database $DATABASE_URL up

# Rollback migration
migrate -path migrations -database $DATABASE_URL down 1
```

## API Examples

### Tạo ticket mới
```bash
curl -X POST http://localhost:8080/api/v2/kanban/tickets \
  -H "Content-Type: application/json" \
  -d '{
    "column_id": "uuid-column-id",
    "title": "New task",
    "description": "Task description",
    "priority": "HIGH"
  }'
```

### Di chuyển ticket với position
```bash
curl -X PATCH http://localhost:8080/api/v2/kanban/tickets/move-with-position \
  -H "Content-Type: application/json" \
  -d '{
    "ticket_id": "uuid-ticket-id",
    "column_id": "uuid-column-id",
    "position": 2
  }'
```

### Tạo assignment
```bash
curl -X POST http://localhost:8080/api/v2/kanban/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "ticket_id": "uuid-ticket-id",
    "user_id": "uuid-user-id",
    "user_full_name": "Nguyễn Văn A"
  }'
```

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

Nếu có vấn đề hoặc câu hỏi, vui lòng tạo issue trong repository này.

---

**Phát triển bởi Gau Team** 🚀

