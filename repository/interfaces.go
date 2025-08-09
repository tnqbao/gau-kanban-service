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
	Create(ticket *entity.Ticket) error
	GetAll() ([]entity.Ticket, error)
	GetByID(id string) (*entity.Ticket, error)
	GetByColumnID(columnID string) ([]entity.Ticket, error)
	Update(ticket *entity.Ticket) error
	Delete(id string) error
	MoveToColumn(ticketID, columnID string) error
	GetWithAssignments(ticketID string) (*entity.Ticket, []entity.TaskAssignment, error)
	GetWithLabels(ticketID string) (*entity.Ticket, []entity.Label, error)
	Search(query string) ([]entity.Ticket, error)
	GetTicketDetail(ticketID string) (*TicketDetailDTO, error)
	GetTicketWithAllRelations(ticketID string) (*TicketDTO, error)
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

// TaskAssignmentRepositoryInterface defines methods for task assignment operations
type TaskAssignmentRepositoryInterface interface {
	Create(assignment *entity.TaskAssignment) error
	GetAll() ([]entity.TaskAssignment, error)
	GetByID(id string) (*entity.TaskAssignment, error)
	GetByTicketID(ticketID string) ([]entity.TaskAssignment, error)
	GetByUserID(userID string) ([]entity.TaskAssignment, error)
	Update(assignment *entity.TaskAssignment) error
	Delete(id string) error
	DeleteByTicketAndUser(ticketID, userID string) error
}

// TicketCommentRepositoryInterface defines methods for ticket comment operations
type TicketCommentRepositoryInterface interface {
	Create(comment *entity.TicketComment) error
	GetAll() ([]entity.TicketComment, error)
	GetByID(id string) (*entity.TicketComment, error)
	GetByTicketID(ticketID string) ([]entity.TicketComment, error)
	GetByUserID(userID string) ([]entity.TicketComment, error)
	Update(comment *entity.TicketComment) error
	Delete(id string) error
}

// TicketLabelRepositoryInterface defines methods for ticket-label relationship operations
type TicketLabelRepositoryInterface interface {
	AddLabelToTicket(ticketID, labelID string) error
	RemoveLabelFromTicket(ticketID, labelID string) error
	GetTicketsByLabelID(labelID string) ([]entity.Ticket, error)
	GetLabelsByTicketID(ticketID string) ([]entity.Label, error)
	RemoveAllLabelsFromTicket(ticketID string) error
}
