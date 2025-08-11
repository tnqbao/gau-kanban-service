package controller

// Checklist DTOs
// Request/Response structures
type KanbanColumnResponse struct {
	ID      string                 `json:"id"`
	Title   string                 `json:"title"`
	Order   int                    `json:"order"`
	Tickets []KanbanTicketResponse `json:"tickets"`
}

type KanbanTicketResponse struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	TicketNo    string   `json:"ticketNo"`
	Tags        []string `json:"tags"`
	Assignees   []string `json:"assignees"`
	Completed   bool     `json:"completed"`
	DueDate     string   `json:"due_date,omitempty"`
	Priority    string   `json:"priority,omitempty"`
}

type CreateColumnRequest struct {
	Name     string `json:"name" binding:"required"`
	Position int    `json:"position"`
}

type UpdateColumnRequest struct {
	Name     string `json:"name"`
	Position *int   `json:"position"`
}

type UpdateColumnPositionRequest struct {
	Position int `json:"position" binding:"required"`
}

type UpdateTicketPositionRequest struct {
	Position int `json:"position" binding:"required"`
}

type CreateAssignmentRequest struct {
	TicketID     string `json:"ticket_id" binding:"required"`
	UserID       string `json:"user_id" binding:"required"`
	UserFullName string `json:"user_full_name" binding:"required"`
}

type UpdateAssignmentRequest struct {
	UserFullName string `json:"user_full_name"`
}

type ChecklistDTO struct {
	ID        string `json:"id"`
	TicketID  string `json:"ticket_id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateChecklistRequest struct {
	TicketID string `json:"ticket_id" validate:"required"`
	Title    string `json:"title" validate:"required"`
}

type UpdateChecklistRequest struct {
	Title     *string `json:"title,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
	Position  *int    `json:"position,omitempty"`
}

type UpdateChecklistPositionRequest struct {
	Position int `json:"position" validate:"required"`
}

// Ticket DTOs with Checklist support
type CreateTicketRequest struct {
	ColumnID    string                   `json:"column_id" validate:"required"`
	Title       string                   `json:"title" validate:"required"`
	Description string                   `json:"description"`
	DueDate     string                   `json:"due_date"`
	Priority    string                   `json:"priority"`
	Checklists  []CreateChecklistRequest `json:"checklists,omitempty"`
}

type UpdateTicketRequest struct {
	Title       *string                  `json:"title,omitempty"`
	Description *string                  `json:"description,omitempty"`
	DueDate     *string                  `json:"due_date,omitempty"`
	Priority    *string                  `json:"priority,omitempty"`
	Checklists  []UpdateChecklistRequest `json:"checklists,omitempty"`
}

type TicketWithChecklistsDTO struct {
	ID          string         `json:"id"`
	TicketNo    string         `json:"ticket_no"`
	ColumnID    string         `json:"column_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	DueDate     string         `json:"due_date"`
	Priority    string         `json:"priority"`
	Position    int            `json:"position"`
	Checklists  []ChecklistDTO `json:"checklists"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

// Existing DTOs
type MoveTicketRequest struct {
	TicketID string `json:"ticket_id" validate:"required"`
	ColumnID string `json:"column_id" validate:"required"`
}
