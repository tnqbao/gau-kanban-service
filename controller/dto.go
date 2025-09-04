package controller

// Board DTOs
type CreateBoardRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	FullName    string `json:"fullname" binding:"required"`
}

type UpdateBoardRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Member DTOs
type AddMemberRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	BoardID string `json:"board_id" binding:"required"`
}

// Column DTOs
type CreateColumnRequest struct {
	BoardID  string `json:"board_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	WipLimit *int   `json:"wip_limit"`
}

type UpdateColumnRequest struct {
	Title    string `json:"title"`
	WipLimit *int   `json:"wip_limit"`
}

// Reorder column request - Optimal fractional approach
type ReorderColumnRequest struct {
	Position string `json:"position" binding:"required"` // "first", "last", "after:column_id", "before:column_id"
}

// Alternative: Batch reorder (if you prefer full control)
type BatchReorderColumnsRequest struct {
	ColumnOrders []ColumnOrder `json:"column_orders" binding:"required"`
}

type ColumnOrder struct {
	ColumnID string `json:"column_id" binding:"required"`
	Order    int    `json:"order" binding:"required"`
}

// Ticket DTOs
type CreateTicketRequest struct {
	ColumnID string `json:"column_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
}

type UpdateTicketRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Priority    *string `json:"priority"`
	DueDate     *string `json:"due_date"` // ISO 8601 format
	Status      *string `json:"status"`
}

type MoveTicketRequest struct {
	ColumnID string `json:"column_id" binding:"required"`
	Position string `json:"position" binding:"required"` // "first", "last", "after:ticket_id", "before:ticket_id"
}

// Checklist DTOs
type CreateChecklistRequest struct {
	TicketID string `json:"ticket_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
}

type UpdateChecklistRequest struct {
	Status string `json:"status" binding:"required"` // "pending", "completed"
}

// Assignee DTOs
type CreateAssigneeRequest struct {
	TicketID string `json:"ticket_id" binding:"required"`
	MemberID string `json:"member_id" binding:"required"`
}

// Label DTOs
type CreateLabelRequest struct {
	BoardID string `json:"board_id" binding:"required"`
	Title   string `json:"title" binding:"required"`
}

// Label Ticket DTOs
type CreateLabelTicketRequest struct {
	LabelID  string `json:"label_id" binding:"required"`
	TicketID string `json:"ticket_id" binding:"required"`
}

// Search and Filter DTOs
type SearchTicketsRequest struct {
	Search string `form:"search"`
}

type FilterTicketsRequest struct {
	Assignee string `form:"assignee"` // member_id
	Label    string `form:"label"`    // label_id
	Status   string `form:"status"`   // ticket status
}

// User DTOs
type CreateUserRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	FullName string `json:"fullname" binding:"required"`
}
