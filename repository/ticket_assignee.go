package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
)

// TicketAssignee methods
func (r *Repository) CreateTicketAssignee(assignee *entity.TicketAssignee) error {
	return r.db.Create(assignee).Error
}

func (r *Repository) GetTicketAssigneesByTicketID(ticketID string) ([]entity.TicketAssignee, error) {
	var assignees []entity.TicketAssignee
	err := r.db.Where("ticket_id = ?", ticketID).
		Preload("Member").
		Find(&assignees).Error
	return assignees, err
}

func (r *Repository) CheckTicketAssigneeExists(ticketID, memberID string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.TicketAssignee{}).
		Where("ticket_id = ? AND member_id = ?", ticketID, memberID).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) DeleteTicketAssignee(ticketID, memberID string) error {
	return r.db.Where("ticket_id = ? AND member_id = ?", ticketID, memberID).
		Delete(&entity.TicketAssignee{}).Error
}
