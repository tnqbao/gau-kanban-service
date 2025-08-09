package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
	"gorm.io/gorm"
)

type TaskAssignmentRepository struct {
	db *gorm.DB
}

func NewTaskAssignmentRepository(db *gorm.DB) TaskAssignmentRepositoryInterface {
	return &TaskAssignmentRepository{db: db}
}

func (r *TaskAssignmentRepository) Create(assignment *entity.TaskAssignment) error {
	return r.db.Create(assignment).Error
}

func (r *TaskAssignmentRepository) GetAll() ([]entity.TaskAssignment, error) {
	var assignments []entity.TaskAssignment
	err := r.db.Order("assigned_at DESC").Find(&assignments).Error
	return assignments, err
}

func (r *TaskAssignmentRepository) GetByID(id string) (*entity.TaskAssignment, error) {
	var assignment entity.TaskAssignment
	err := r.db.Where("id = ?", id).First(&assignment).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

func (r *TaskAssignmentRepository) GetByTicketID(ticketID string) ([]entity.TaskAssignment, error) {
	var assignments []entity.TaskAssignment
	err := r.db.Where("ticket_id = ?", ticketID).Order("assigned_at ASC").Find(&assignments).Error
	return assignments, err
}

func (r *TaskAssignmentRepository) GetByUserID(userID string) ([]entity.TaskAssignment, error) {
	var assignments []entity.TaskAssignment
	err := r.db.Where("user_id = ?", userID).Order("assigned_at DESC").Find(&assignments).Error
	return assignments, err
}

func (r *TaskAssignmentRepository) Update(assignment *entity.TaskAssignment) error {
	return r.db.Save(assignment).Error
}

func (r *TaskAssignmentRepository) Delete(id string) error {
	return r.db.Delete(&entity.TaskAssignment{}, "id = ?", id).Error
}

func (r *TaskAssignmentRepository) DeleteByTicketAndUser(ticketID, userID string) error {
	return r.db.Where("ticket_id = ? AND user_id = ?", ticketID, userID).Delete(&entity.TaskAssignment{}).Error
}
