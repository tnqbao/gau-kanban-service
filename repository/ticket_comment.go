package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
)

func (r *Repository) CreateComment(comment *entity.TicketComment) error {
	return r.db.Create(comment).Error
}

func (r *Repository) GetAllComment() ([]entity.TicketComment, error) {
	var comments []entity.TicketComment
	err := r.db.Order("created_at DESC").Find(&comments).Error
	return comments, err
}

func (r *Repository) GetCommentByID(id string) (*entity.TicketComment, error) {
	var comment entity.TicketComment
	err := r.db.Where("id = ?", id).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *Repository) GetCommentByTicketID(ticketID string) ([]entity.TicketComment, error) {
	var comments []entity.TicketComment
	err := r.db.Where("ticket_id = ?", ticketID).Order("created_at ASC").Find(&comments).Error
	return comments, err
}

func (r *Repository) GetCommentByUserID(userID string) ([]entity.TicketComment, error) {
	var comments []entity.TicketComment
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&comments).Error
	return comments, err
}

func (r *Repository) UpdateComment(comment *entity.TicketComment) error {
	return r.db.Save(comment).Error
}

func (r *Repository) DeleteComment(id string) error {
	return r.db.Delete(&entity.TicketComment{}, "id = ?", id).Error
}
