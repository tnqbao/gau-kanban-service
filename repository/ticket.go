package repository

import (
	"fmt"

	"github.com/tnqbao/gau-kanban-service/entity"
)

// Ticket methods
func (r *Repository) CreateTicket(ticket *entity.Ticket) error {
	// Generate ticket number
	var count int64
	r.db.Model(&entity.Ticket{}).Count(&count)
	ticket.TicketNumber = fmt.Sprintf("#%06d", count+1)

	// Initialize empty arrays
	if ticket.Assignees == nil {
		ticket.Assignees = []entity.TicketAssignee{}
	}
	if ticket.Labels == nil {
		ticket.Labels = []entity.TicketLabel{}
	}
	if ticket.Checklists == nil {
		ticket.Checklists = []entity.Checklist{}
	}

	return r.db.Create(ticket).Error
}

func (r *Repository) GenerateNextTicketNumber() (string, error) {
	var count int64
	err := r.db.Model(&entity.Ticket{}).Count(&count).Error
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("#%06d", count+1), nil
}

func (r *Repository) GetMaxOrderByColumnID(columnID string) (int, error) {
	var maxPosition int
	err := r.db.Model(&entity.Ticket{}).
		Where("column_id = ?", columnID).
		Select("COALESCE(MAX(position), 0)").
		Scan(&maxPosition).Error
	return maxPosition, err
}

func (r *Repository) GetTicketByID(id string) (*entity.Ticket, error) {
	var ticket entity.Ticket
	err := r.db.Where("id = ?", id).
		Preload("Assignees.Member.User").
		Preload("Labels.Label").
		Preload("Checklists").
		First(&ticket).Error
	if err != nil {
		return nil, err
	}

	// Ensure arrays are never nil
	if ticket.Assignees == nil {
		ticket.Assignees = []entity.TicketAssignee{}
	}
	if ticket.Labels == nil {
		ticket.Labels = []entity.TicketLabel{}
	}
	if ticket.Checklists == nil {
		ticket.Checklists = []entity.Checklist{}
	}

	return &ticket, nil
}

func (r *Repository) GetTicketsByColumnID(columnID string) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	err := r.db.Where("column_id = ?", columnID).
		Order("position ASC").
		Preload("Assignees.Member.User").
		Preload("Labels.Label").
		Preload("Checklists").
		Find(&tickets).Error

	// Ensure arrays are never nil for each ticket
	for i := range tickets {
		if tickets[i].Assignees == nil {
			tickets[i].Assignees = []entity.TicketAssignee{}
		}
		if tickets[i].Labels == nil {
			tickets[i].Labels = []entity.TicketLabel{}
		}
		if tickets[i].Checklists == nil {
			tickets[i].Checklists = []entity.Checklist{}
		}
	}

	return tickets, err
}

func (r *Repository) UpdateTicket(ticket *entity.Ticket) error {
	return r.db.Save(ticket).Error
}

func (r *Repository) DeleteTicket(id string) error {
	return r.db.Delete(&entity.Ticket{}, "id = ?", id).Error
}

func (r *Repository) GetNextTicketPosition(columnID string) (int, error) {
	var maxPosition int
	err := r.db.Model(&entity.Ticket{}).
		Where("column_id = ?", columnID).
		Select("COALESCE(MAX(position), 0) + 1").
		Scan(&maxPosition).Error
	return maxPosition, err
}

func (r *Repository) GetMinPositionByColumnID(columnID string) (int, error) {
	var minPosition int
	err := r.db.Model(&entity.Ticket{}).
		Where("column_id = ?", columnID).
		Select("COALESCE(MIN(position), 0)").
		Scan(&minPosition).Error
	return minPosition, err
}

func (r *Repository) GetMaxPositionByColumnID(columnID string) (int, error) {
	var maxPosition int
	err := r.db.Model(&entity.Ticket{}).
		Where("column_id = ?", columnID).
		Select("COALESCE(MAX(position), 0)").
		Scan(&maxPosition).Error
	return maxPosition, err
}

func (r *Repository) GetNextPositionAfter(columnID string, position int) (int, error) {
	var nextPosition int
	err := r.db.Model(&entity.Ticket{}).
		Where("column_id = ? AND position > ?", columnID, position).
		Select("MIN(position)").
		Scan(&nextPosition).Error
	return nextPosition, err
}

func (r *Repository) GetPreviousPositionBefore(columnID string, position int) (int, error) {
	var prevPosition int
	err := r.db.Model(&entity.Ticket{}).
		Where("column_id = ? AND position < ?", columnID, position).
		Select("MAX(position)").
		Scan(&prevPosition).Error
	return prevPosition, err
}

func (r *Repository) SearchTicketsByColumnAndTitle(columnID string, searchTerm string) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	query := r.db.Where("column_id = ?", columnID)

	if searchTerm != "" {
		// Search in title and ticket_number (case-insensitive)
		query = query.Where("LOWER(title) LIKE LOWER(?) OR LOWER(ticket_number) LIKE LOWER(?)",
			"%"+searchTerm+"%", "%"+searchTerm+"%")
	}

	err := query.Order("position ASC").
		Preload("Assignees.Member.User").
		Preload("Labels.Label").
		Preload("Checklists").
		Find(&tickets).Error
	return tickets, err
}

func (r *Repository) FilterTicketsByColumn(columnID string, assigneeID string, labelID string, status string) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	query := r.db.Where("column_id = ?", columnID)

	// Filter by status if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// If filtering by assignee, we need to join with ticket_assignees
	if assigneeID != "" {
		query = query.Joins("INNER JOIN ticket_assignees ON tickets.id = ticket_assignees.ticket_id").
			Where("ticket_assignees.member_id = ?", assigneeID)
	}

	// If filtering by label, we need to join with ticket_labels
	if labelID != "" {
		if assigneeID != "" {
			// Already have a join, add another condition
			query = query.Joins("INNER JOIN ticket_labels ON tickets.id = ticket_labels.ticket_id").
				Where("ticket_labels.label_id = ?", labelID)
		} else {
			// First join
			query = query.Joins("INNER JOIN ticket_labels ON tickets.id = ticket_labels.ticket_id").
				Where("ticket_labels.label_id = ?", labelID)
		}
	}

	err := query.Order("position ASC").
		Preload("Assignees.Member.User").
		Preload("Labels.Label").
		Preload("Checklists").
		Find(&tickets).Error
	return tickets, err
}
