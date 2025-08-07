package repository

// DTO structures for API responses

// TicketDTO với đầy đủ thông tin từ các bảng liên quan
type TicketDTO struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	TicketID    string        `json:"ticket_id"`
	Labels      []LabelDTO    `json:"labels"`    // Thông tin chi tiết labels
	Assignees   []AssigneeDTO `json:"assignees"` // Thông tin chi tiết assignees
	Comments    []CommentDTO  `json:"comments"`  // Danh sách comments
	Completed   bool          `json:"completed"`
	DueDate     string        `json:"due_date,omitempty"`
	Priority    string        `json:"priority,omitempty"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
}

// ColumnWithTicketsDTO cho kanban board
type ColumnWithTicketsDTO struct {
	ID      string      `json:"id"`
	Title   string      `json:"title"`
	Order   int         `json:"order"`
	Tickets []TicketDTO `json:"tickets"`
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
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
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
