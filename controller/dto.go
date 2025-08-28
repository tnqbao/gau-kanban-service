package controller

type KanbanColumnResponse struct {
	ID       string                 `json:"id"`
	Title    string                 `json:"title"`
	Position int                    `json:"position"`
	Tickets  []KanbanTicketResponse `json:"tickets"`
}

type KanbanTicketResponse struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	TicketNo    string   `json:"ticketNo"`
	Labels      []string `json:"Label"`
	Assignees   []string `json:"assignees"`
	Completed   bool     `json:"completed"`
	DueDate     *string  `json:"due_date,omitempty"`
	Priority    *string  `json:"priority,omitempty"`
}

type CreateColumnRequest struct {
	Title   string `json:"title" binding:"required"`
	BoardID string `json:"board_id" binding:"required"`
}

type UpdateColumnRequest struct {
	Title    string `json:"title"`
	Position *int   `json:"position"`
}

type UpdateColumnPositionRequest struct {
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

// Ticket DTOs
type CreateTicketRequest struct {
	ColumnID string `json:"column_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	//Description string                     `json:"description"`
	//DueDate     string                     `json:"due_date"`
	//Priority    string                     `json:"priority"`
	//Assignments []CreateAssignmentInTicket `json:"assignments"`
	//Checklists  []CreateChecklistInTicket  `json:"checklists"`
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
	DueDate     *string         `json:"due_date"`
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

// Advanced position change requests
type ChangeColumnPositionRequest struct {
	NewPosition int `json:"new_position" binding:"required"`
}

type ChangeTicketPositionRequest struct {
	NewColumnID string `json:"new_column_id" binding:"required"`
	NewPosition int    `json:"new_position" binding:"required"`
}

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
type CreateMemberRequest struct {
	BoardID  string `json:"board_id" binding:"required"`
	MemberID string `json:"member_id" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
}

type UpdateMemberRequest struct {
	FullName string `json:"full_name"`
}

// BoardWithMembers response DTO
type BoardWithMembersResponse struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Archived    bool        `json:"archived"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	Members     []MemberDTO `json:"members"`
}

type MemberDTO struct {
	ID        string `json:"id"`
	BoardID   string `json:"board_id"`
	MemberID  string `json:"member_id"`
	FullName  string `json:"full_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
