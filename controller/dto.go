package controller

// Ticket DTOs
type CreateTicketRequest struct {
	ColumnID    string                     `json:"column_id" binding:"required"`
	Title       string                     `json:"title" binding:"required"`
	Description string                     `json:"description"`
	DueDate     string                     `json:"due_date"`
	Priority    string                     `json:"priority"`
	Assignments []CreateAssignmentInTicket `json:"assignments"`
	Checklists  []CreateChecklistInTicket  `json:"checklists"`
}

type UpdateTicketRequest struct {
	Title       *string                    `json:"title"`
	Description *string                    `json:"description"`
	DueDate     *string                    `json:"due_date"`
	Priority    *string                    `json:"priority"`
	Assignments []CreateAssignmentInTicket `json:"assignments"`
	Checklists  []CreateChecklistInTicket  `json:"checklists"`
}

type CreateAssignmentInTicket struct {
	UserID       string `json:"user_id" binding:"required"`
	UserFullName string `json:"user_full_name" binding:"required"`
}

type CreateChecklistInTicket struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type UpdateTicketPositionRequest struct {
	ColumnID string `json:"column_id" binding:"required"`
	Position int    `json:"position" binding:"required"`
}

type MoveTicketRequest struct {
	TicketID string `json:"ticket_id" binding:"required"`
	ColumnID string `json:"column_id" binding:"required"`
}

type MoveTicketWithPositionRequest struct {
	TicketID string `json:"ticket_id" binding:"required"`
	ColumnID string `json:"column_id" binding:"required"`
	Position int    `json:"position" binding:"required"`
}

// Checklist DTOs
type CreateChecklistRequest struct {
	TicketID string `json:"ticket_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
}

type UpdateChecklistRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
	Position  *int    `json:"position"`
}

type UpdateChecklistPositionRequest struct {
	Position int `json:"position" binding:"required"`
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

// Response DTOs with assignments and checklists
type TicketWithDetailsResponse struct {
	ID          string          `json:"id"`
	TicketNo    string          `json:"ticket_no"`
	ColumnID    string          `json:"column_id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	DueDate     string          `json:"due_date"`
	Priority    string          `json:"priority"`
	Position    int             `json:"position"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
	Assignments []AssignmentDTO `json:"assignments"`
	Checklists  []ChecklistDTO  `json:"checklists"`
}

type AssignmentDTO struct {
	ID           string `json:"id"`
	TicketID     string `json:"ticket_id"`
	UserID       string `json:"user_id"`
	UserFullName string `json:"user_full_name"`
	AssignedAt   string `json:"assigned_at"`
}
