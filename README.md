t # Gau Kanban Service

## Mô tả dự án

Gau Kanban Service là một RESTful API service được xây dựng bằng Go để quản lý bảng Kanban. Dự án cung cấp các tính năng hoàn chỉnh để tạo và quản lý boards, columns, tickets, assignments, checklists và labels cho việc quản lý dự án theo phương pháp Kanban.

## Tính năng chính

### 🏗️ Quản lý Columns
- Tạo, sửa, xóa columns
- Sắp xếp lại vị trí columns
- Quản lý thứ tự hiển thị

### 🎫 Quản lý Tickets
- Tạo tickets với ticket number tự động (TASK-XXXX format)
- CRUD operations cho tickets
- Di chuyển tickets giữa các columns với position management thông minh
- Drag & drop hỗ trợ di chuyển vào vị trí bất kỳ trong column
- Tự động sắp xếp vị trí khi tạo ticket mới (luôn ở cuối column)
- Hỗ trợ due date và priority
- Tích hợp assignments và checklists trong ticket operations

### 👥 Quản lý Assignments
- Gán người dùng vào tickets (không cần tạo bảng user riêng)
- Quản lý thông tin assignees (user_id, user_full_name)
- CRUD operations: tạo, sửa, xóa assignments
- Xóa tất cả assignments theo user ID
- Hiển thị assignees trong thông tin tickets

### ✅ Quản lý Checklists
- Tạo checklist items cho tickets
- Đánh dấu hoàn thành/chưa hoàn thành
- Sắp xếp thứ tự checklist items
- Tích hợp trong ticket create/update operations
- CRUD operations riêng biệt cho từng checklist item

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
├── controller/            # HTTP controllers
│   ├── ticket.go         # Ticket operations
│   ├── assignment.go     # Assignment operations  
│   ├── checklist.go      # Checklist operations
│   ├── column.go         # Column operations
│   └── dto.go            # Data Transfer Objects
├── entity/               # Domain entities
│   ├── ticket.go         # Ticket entity
│   ├── task_assignment.go # Assignment entity
│   ├── checklist.go      # Checklist entity
│   └── column.go         # Column entity
├── repository/           # Data access layer
│   ├── ticket.go         # Ticket repository
│   ├── task_assignment.go # Assignment repository
│   ├── checklist.go      # Checklist repository
│   └── interfaces.go     # Repository interfaces
├── routes/               # Route definitions
├── migrations/           # Database migrations
├── utils/                # Utility functions
└── deploy/               # Kubernetes deployment configs
```

## API Endpoints

### Tickets
- `POST /api/tickets` - Tạo ticket mới (có thể kèm assignments và checklists)
- `GET /api/tickets` - Lấy danh sách tickets (kèm assignees và checklists)
- `GET /api/tickets/:id` - Lấy ticket theo ID (kèm assignees và checklists)
- `PUT /api/tickets/:id` - Cập nhật ticket (có thể cập nhật assignments và checklists)
- `DELETE /api/tickets/:id` - Xóa ticket
- `PUT /api/tickets/:id/position` - Cập nhật vị trí ticket trong column
- `PUT /api/tickets/move` - Di chuyển ticket sang column khác
- `PUT /api/tickets/move-with-position` - Di chuyển ticket với vị trí cụ thể

### Assignments
- `POST /api/assignments` - Tạo assignment mới
- `GET /api/assignments/ticket/:ticket_id` - Lấy assignments của ticket
- `PUT /api/assignments/:id` - Cập nhật assignment
- `DELETE /api/assignments/:id` - Xóa assignment
- `DELETE /api/assignments/user/:user_id` - Xóa tất cả assignments của user

### Checklists
- `POST /api/checklists` - Tạo checklist item mới
- `GET /api/checklists/ticket/:ticketId` - Lấy checklists của ticket
- `PUT /api/checklists/:id` - Cập nhật checklist item
- `PUT /api/checklists/:id/position` - Cập nhật vị trí checklist item
- `DELETE /api/checklists/:id` - Xóa checklist item

### Columns
- `POST /api/columns` - Tạo column mới
- `GET /api/columns` - Lấy danh sách columns
- `GET /api/columns/:id` - Lấy column theo ID
- `PUT /api/columns/:id` - Cập nhật column
- `DELETE /api/columns/:id` - Xóa column
- `PUT /api/columns/:id/position` - Cập nhật vị trí column

### Kanban Board
- `GET /api/kanban/board` - Lấy toàn bộ kanban board với columns và tickets

## Tính năng đặc biệt

### Smart Position Management
- Hệ thống tự động quản lý vị trí tickets khi drag & drop
- Hỗ trợ di chuyển ticket vào vị trí bất kỳ trong column (ví dụ: từ vị trí 2 lên vị trí 4 trong column có 10 tickets)
- Tự động điều chỉnh position của các tickets khác
- Xử lý di chuyển giữa các columns khác nhau

### Automatic Ticket Numbering
- Ticket number tự động theo format TASK-XXXX (ví dụ: TASK-0001, TASK-0002)
- Sử dụng PostgreSQL sequence để đảm bảo tính duy nhất
- Không bị trùng lặp khi tạo đồng thời

### Integrated Operations
- Tạo/cập nhật ticket có thể kèm theo assignments và checklists
- Tự động xóa assignments và checklists khi xóa ticket
- API riêng biệt cho từng component để tối ưu performance

## Cài đặt và chạy

### Prerequisites
- Go 1.21+
- PostgreSQL 13+
- Docker & Docker Compose (optional)

### Local Development
```bash
# Clone repository
git clone <repository-url>
cd gau-kanban-service

# Install dependencies
go mod tidy

# Setup database
createdb kanban_db

# Run migrations
migrate -path migrations -database "postgres://username:password@localhost/kanban_db?sslmode=disable" up

# Run the application
go run main.go
```

### Docker Development
```bash
# Build and run with Docker Compose
docker-compose up --build

# Run migrations
docker-compose exec app migrate -path migrations -database $DATABASE_URL up
```

### Kubernetes Deployment
```bash
# Deploy to staging
cd deploy/k8s/staging
./apply.sh

# Deploy to production
cd deploy/k8s/production
./apply.sh
```

## Database Schema

### Tickets Table
- id (UUID, Primary Key)
- ticket_no (Text, Unique) - Format: TASK-XXXX
- column_id (UUID, Foreign Key)
- title (Text)
- description (Text)
- due_date (Date)
- priority (Text)
- position (Integer) - Vị trí trong column
- created_at, updated_at (Timestamp)

### Task Assignments Table
- id (UUID, Primary Key)
- ticket_id (UUID, Foreign Key)
- user_id (UUID)
- user_full_name (Text)
- assigned_at (Timestamp)

### Checklists Table
- id (UUID, Primary Key)
- ticket_id (UUID, Foreign Key)
- title (Text)
- completed (Boolean)
- position (Integer) - Vị trí trong ticket
- created_at, updated_at (Timestamp)

## Environment Variables

```bash
DATABASE_URL=postgres://username:password@localhost/kanban_db?sslmode=disable
PORT=8080
GIN_MODE=release # for production
```

## Contributing

1. Fork repository
2. Tạo feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Tạo Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
