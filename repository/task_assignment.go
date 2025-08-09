package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
)

func (r *Repository) CreateAssignment(assignment *entity.TaskAssignment) error {
	return r.db.Create(assignment).Error
}

func (r *Repository) GetAllAssignment() ([]entity.TaskAssignment, error) {
	var assignments []entity.TaskAssignment
	err := r.db.Order("assigned_at DESC").Find(&assignments).Error
	return assignments, err
}

func (r *Repository) GetAssignmentByID(id string) (*entity.TaskAssignment, error) {
	var assignment entity.TaskAssignment
	err := r.db.Where("id = ?", id).First(&assignment).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

func (r *Repository) GetAssignmentByTicketID(ticketID string) ([]entity.TaskAssignment, error) {
	var assignments []entity.TaskAssignment
	err := r.db.Where("ticket_id = ?", ticketID).Order("assigned_at ASC").Find(&assignments).Error
	return assignments, err
}

func (r *Repository) GetAssignmentsByTicketID(ticketID string) ([]entity.TaskAssignment, error) {
	return r.GetAssignmentByTicketID(ticketID)
}

func (r *Repository) GetAssignmentByUserID(userID string) ([]entity.TaskAssignment, error) {
	var assignments []entity.TaskAssignment
	err := r.db.Where("user_id = ?", userID).Order("assigned_at DESC").Find(&assignments).Error
	return assignments, err
}

func (r *Repository) UpdateAssignment(assignment *entity.TaskAssignment) error {
	return r.db.Save(assignment).Error
}

func (r *Repository) DeleteAssignment(id string) error {
	return r.db.Delete(&entity.TaskAssignment{}, "id = ?", id).Error
}

func (r *Repository) DeleteAssignmentByTicketAndUser(ticketID, userID string) error {
	return r.db.Where("ticket_id = ? AND user_id = ?", ticketID, userID).Delete(&entity.TaskAssignment{}).Error
}

func (r *Repository) DeleteAssignmentsByUserID(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&entity.TaskAssignment{}).Error
}
