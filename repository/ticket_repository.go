package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepositoryInterface {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Create(ticket *entity.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *TicketRepository) GetAll() ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	err := r.db.Order("created_at DESC").Find(&tickets).Error
	return tickets, err
}

func (r *TicketRepository) GetByID(id string) (*entity.Ticket, error) {
	var ticket entity.Ticket
	err := r.db.Where("id = ?", id).First(&ticket).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepository) GetByColumnID(columnID string) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	err := r.db.Where("column_id = ?", columnID).Order("created_at ASC").Find(&tickets).Error
	return tickets, err
}

func (r *TicketRepository) Update(ticket *entity.Ticket) error {
	return r.db.Save(ticket).Error
}

func (r *TicketRepository) Delete(id string) error {
	return r.db.Delete(&entity.Ticket{}, "id = ?", id).Error
}

func (r *TicketRepository) MoveToColumn(ticketID, columnID string) error {
	return r.db.Model(&entity.Ticket{}).Where("id = ?", ticketID).Update("column_id", columnID).Error
}

func (r *TicketRepository) GetWithAssignments(ticketID string) (*entity.Ticket, []entity.TaskAssignment, error) {
	var ticket entity.Ticket
	var assignments []entity.TaskAssignment

	err := r.db.Where("id = ?", ticketID).First(&ticket).Error
	if err != nil {
		return nil, nil, err
	}

	err = r.db.Where("ticket_id = ?", ticketID).Find(&assignments).Error
	if err != nil {
		return &ticket, nil, err
	}

	return &ticket, assignments, nil
}

func (r *TicketRepository) GetWithLabels(ticketID string) (*entity.Ticket, []entity.Label, error) {
	var ticket entity.Ticket
	var labels []entity.Label

	err := r.db.Where("id = ?", ticketID).First(&ticket).Error
	if err != nil {
		return nil, nil, err
	}

	err = r.db.Table("labels").
		Joins("JOIN ticket_labels ON labels.id = ticket_labels.label_id").
		Where("ticket_labels.ticket_id = ?", ticketID).
		Find(&labels).Error

	if err != nil {
		return &ticket, nil, err
	}

	return &ticket, labels, nil
}

func (r *TicketRepository) Search(query string) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	searchPattern := "%" + query + "%"
	err := r.db.Where("title ILIKE ? OR description ILIKE ?", searchPattern, searchPattern).
		Order("created_at DESC").Find(&tickets).Error
	return tickets, err
}

func (r *TicketRepository) GetTicketDetail(ticketID string) (*TicketDetailDTO, error) {
	var ticket entity.Ticket
	var column entity.Column

	// Lấy thông tin ticket
	err := r.db.Where("id = ?", ticketID).First(&ticket).Error
	if err != nil {
		return nil, err
	}

	// Lấy thông tin column
	err = r.db.Where("id = ?", ticket.ColumnID).First(&column).Error
	if err != nil {
		return nil, err
	}

	// Lấy labels chi tiết
	var labels []LabelDTO
	err = r.db.Table("labels").
		Select("labels.id, labels.name, labels.color").
		Joins("JOIN ticket_labels ON labels.id = ticket_labels.label_id").
		Where("ticket_labels.ticket_id = ?", ticket.ID).
		Scan(&labels).Error
	if err != nil {
		labels = []LabelDTO{}
	}

	// Lấy assignees chi tiết
	var assignees []AssigneeDTO
	err = r.db.Table("task_assignments").
		Select("task_assignments.id, task_assignments.user_id, task_assignments.user_full_name as user_full_name, task_assignments.assigned_at").
		Where("ticket_id = ?", ticket.ID).
		Scan(&assignees).Error
	if err != nil {
		assignees = []AssigneeDTO{}
	}

	// Lấy comments chi tiết
	var comments []CommentDTO
	err = r.db.Table("ticket_comments").
		Select("id, user_id, content, created_at").
		Where("ticket_id = ?", ticket.ID).
		Order("created_at ASC").
		Scan(&comments).Error
	if err != nil {
		comments = []CommentDTO{}
	}

	// Xác định completed
	completed := column.Name == "DONE" || column.Name == "COMPLETED"

	ticketDTO := TicketDTO{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		TicketID:    ticket.ID,
		Labels:      labels,
		Assignees:   assignees,
		Comments:    comments,
		Completed:   completed,
		DueDate:     ticket.DueDate,
		Priority:    ticket.Priority,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
	}

	columnDTO := ColumnDTO{
		ID:        column.ID,
		Name:      column.Name,
		Position:  column.Position,
		CreatedAt: column.CreatedAt,
		UpdatedAt: column.UpdatedAt,
	}

	result := &TicketDetailDTO{
		Ticket:    ticketDTO,
		Column:    columnDTO,
		Labels:    labels,
		Assignees: assignees,
		Comments:  comments,
	}

	return result, nil
}

func (r *TicketRepository) GetTicketWithAllRelations(ticketID string) (*TicketDTO, error) {
	var ticket entity.Ticket
	var column entity.Column

	// Lấy thông tin ticket
	err := r.db.Where("id = ?", ticketID).First(&ticket).Error
	if err != nil {
		return nil, err
	}

	// Lấy thông tin column để xác định completed
	err = r.db.Where("id = ?", ticket.ColumnID).First(&column).Error
	if err != nil {
		return nil, err
	}

	// Lấy labels chi tiết
	var labels []LabelDTO
	err = r.db.Table("labels").
		Select("labels.id, labels.name, labels.color").
		Joins("JOIN ticket_labels ON labels.id = ticket_labels.label_id").
		Where("ticket_labels.ticket_id = ?", ticket.ID).
		Scan(&labels).Error
	if err != nil {
		labels = []LabelDTO{}
	}

	// Lấy assignees chi tiết
	var assignees []AssigneeDTO
	err = r.db.Table("task_assignments").
		Select("task_assignments.id, task_assignments.user_id, task_assignments.user_full_name as user_full_name, task_assignments.assigned_at").
		Where("ticket_id = ?", ticket.ID).
		Scan(&assignees).Error
	if err != nil {
		assignees = []AssigneeDTO{}
	}

	// Lấy comments chi tiết
	var comments []CommentDTO
	err = r.db.Table("ticket_comments").
		Select("id, user_id, content, created_at").
		Where("ticket_id = ?", ticket.ID).
		Order("created_at ASC").
		Scan(&comments).Error
	if err != nil {
		comments = []CommentDTO{}
	}

	// Xác định completed
	completed := column.Name == "DONE" || column.Name == "COMPLETED"

	result := &TicketDTO{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		TicketID:    ticket.ID,
		Labels:      labels,
		Assignees:   assignees,
		Comments:    comments,
		Completed:   completed,
		DueDate:     ticket.DueDate,
		Priority:    ticket.Priority,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
	}

	return result, nil
}
