package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
)

// ColumnRepositoryInterface defines methods for column operations
type ColumnRepositoryInterface interface {
	Create(column *entity.Column) error
	GetAll() ([]entity.Column, error)
	GetByID(id string) (*entity.Column, error)
	Update(column *entity.Column) error
	Delete(id string) error
	UpdatePosition(id string, position int) error
	GetAllWithTickets() ([]ColumnWithTicketsDTO, error)
	GetAllWithFullTicketDetails() ([]ColumnWithTicketsDTO, error)
}

// TicketRepositoryInterface defines methods for ticket operations
type TicketRepositoryInterface interface {
	CreateTicket(ticket *entity.Ticket) error
	GetAllTickets() ([]entity.Ticket, error)
	GetTicketByID(id string) (*entity.Ticket, error)
	GetTicketsByColumnID(columnID string) ([]entity.Ticket, error)
	UpdateTicket(ticket *entity.Ticket) error
	DeleteTicket(id string) error
	MoveTicketToColumn(ticketID, columnID string) error
	MoveTicketToColumnWithPosition(ticketID, columnID string, position int) error
	UpdateTicketPosition(ticketID, columnID string, position int) error
	GetMaxTicketPositionInColumn(columnID string) (int, error)
	GenerateTicketNumber() (string, error)
	GetTicketWithDetails(ticketID string) (*TicketWithDetailsResponse, error)
}

// TaskAssignmentRepositoryInterface defines methods for task assignment operations
type TaskAssignmentRepositoryInterface interface {
	CreateAssignment(assignment *entity.TaskAssignment) error
	GetAllAssignments() ([]entity.TaskAssignment, error)
	GetAssignmentByID(id string) (*entity.TaskAssignment, error)
	GetAssignmentsByTicketID(ticketID string) ([]entity.TaskAssignment, error)
	GetAssignmentsByUserID(userID string) ([]entity.TaskAssignment, error)
	UpdateAssignment(assignment *entity.TaskAssignment) error
	DeleteAssignment(id string) error
	DeleteAssignmentsByTicketID(ticketID string) error
	DeleteAssignmentsByUserID(userID string) error
}

// ChecklistRepositoryInterface defines methods for checklist operations
type ChecklistRepositoryInterface interface {
	CreateChecklist(checklist *entity.Checklist) error
	GetChecklistByID(id string) (*entity.Checklist, error)
	GetChecklistsByTicketID(ticketID string) ([]entity.Checklist, error)
	UpdateChecklist(checklist *entity.Checklist) error
	UpdateChecklistPosition(checklistID string, position int) error
	DeleteChecklist(id string) error
	DeleteChecklistsByTicketID(ticketID string) error
	GetMaxChecklistPosition(ticketID string) (int, error)
}

// LabelRepositoryInterface defines methods for label operations
type LabelRepositoryInterface interface {
	Create(label *entity.Label) error
	GetAll() ([]entity.Label, error)
	GetByID(id string) (*entity.Label, error)
	Update(label *entity.Label) error
	Delete(id string) error
	GetByTicketID(ticketID string) ([]entity.Label, error)
}

// Repository interface combining all repository interfaces
type RepositoryInterface interface {
	TicketRepositoryInterface
	TaskAssignmentRepositoryInterface
	ChecklistRepositoryInterface
	LabelRepositoryInterface
}

// Response DTOs
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

type ChecklistDTO struct {
	ID        string `json:"id"`
	TicketID  string `json:"ticket_id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Position  int    `json:"position"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
