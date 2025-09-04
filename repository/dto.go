package repository

// DTO structures for API responses

// TicketDTO với đầy đủ thông tin từ các bảng liên quan
type TicketDTO struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	TicketID    string         `json:"ticket_id"`  // Maps to TicketNo from entity
	ColumnID    string         `json:"column_id"`  // Added missing field
	Position    int            `json:"position"`   // Added missing field
	Labels      []LabelDTO     `json:"labels"`     // Thông tin chi tiết labels
	Assignees   []AssigneeDTO  `json:"assignees"`  // Thông tin chi tiết assignees
	Comments    []CommentDTO   `json:"comments"`   // Danh sách comments
	Checklists  []ChecklistDTO `json:"checklists"` // Danh sách checklists
	Completed   bool           `json:"completed"`
	DueDate     *string        `json:"due_date,omitempty"`
	Priority    *string        `json:"priority,omitempty"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

// ColumnWithTicketsDTO cho kanban board
type ColumnWithTicketsDTO struct {
	ID        string      `json:"id"`
	Title     string      `json:"title"`
	BoardID   string      `json:"board_id"`
	Position  int         `json:"position"`
	Tickets   []TicketDTO `json:"tickets"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

// AssigneeDTO chứa thông tin chi tiết về người được assign
type AssigneeDTO struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	UserFullName string `json:"user_full_name"`
	AssignedAt   string `json:"assigned_at"`
}

// LabelDTO chứa thông tin chi tiết về label
type LabelDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// CommentDTO chứa thông tin về comment
type CommentDTO struct {
	ID        string `json:"id"`
	TicketID  string `json:"ticket_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// ChecklistDTO chứa thông tin về checklist
type ChecklistDTO struct {
	ID        string `json:"id"`
	TicketID  string `json:"ticket_id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// TicketDetailDTO để lấy thông tin chi tiết một ticket
type TicketDetailDTO struct {
	Ticket    TicketDTO     `json:"ticket"`
	Column    ColumnDTO     `json:"column"`
	Labels    []LabelDTO    `json:"labels"`
	Assignees []AssigneeDTO `json:"assignees"`
	Comments  []CommentDTO  `json:"comments"`
}

// ColumnDTO thông tin cơ bản về column
type ColumnDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
