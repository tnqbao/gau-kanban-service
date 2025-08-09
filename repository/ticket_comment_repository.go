package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
	"gorm.io/gorm"
)

type TicketCommentRepository struct {
	db *gorm.DB
}

func NewTicketCommentRepository(db *gorm.DB) TicketCommentRepositoryInterface {
	return &TicketCommentRepository{db: db}
}

func (r *TicketCommentRepository) Create(comment *entity.TicketComment) error {
	return r.db.Create(comment).Error
}

func (r *TicketCommentRepository) GetAll() ([]entity.TicketComment, error) {
	var comments []entity.TicketComment
	err := r.db.Order("created_at DESC").Find(&comments).Error
	return comments, err
}

func (r *TicketCommentRepository) GetByID(id string) (*entity.TicketComment, error) {
	var comment entity.TicketComment
	err := r.db.Where("id = ?", id).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *TicketCommentRepository) GetByTicketID(ticketID string) ([]entity.TicketComment, error) {
	var comments []entity.TicketComment
	err := r.db.Where("ticket_id = ?", ticketID).Order("created_at ASC").Find(&comments).Error
	return comments, err
}

func (r *TicketCommentRepository) GetByUserID(userID string) ([]entity.TicketComment, error) {
	var comments []entity.TicketComment
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&comments).Error
	return comments, err
}

func (r *TicketCommentRepository) Update(comment *entity.TicketComment) error {
	return r.db.Save(comment).Error
}

func (r *TicketCommentRepository) Delete(id string) error {
	return r.db.Delete(&entity.TicketComment{}, "id = ?", id).Error
}
