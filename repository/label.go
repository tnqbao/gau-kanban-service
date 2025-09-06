package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
)

// Label methods
func (r *Repository) CreateLabel(label *entity.Label) error {
	return r.db.Create(label).Error
}

func (r *Repository) GetLabelsByBoardID(boardID string) ([]entity.Label, error) {
	var labels []entity.Label
	err := r.db.Where("board_id = ?", boardID).Find(&labels).Error
	return labels, err
}

func (r *Repository) GetLabelByID(id string) (*entity.Label, error) {
	var label entity.Label
	err := r.db.Where("id = ?", id).First(&label).Error
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (r *Repository) UpdateLabel(label *entity.Label) error {
	return r.db.Save(label).Error
}

func (r *Repository) DeleteLabel(id string) error {
	return r.db.Delete(&entity.Label{}, "id = ?", id).Error
}

// TicketLabel methods
func (r *Repository) CreateTicketLabel(ticketLabel *entity.TicketLabel) error {
	return r.db.Create(ticketLabel).Error
}

func (r *Repository) GetTicketLabelsByTicketID(ticketID string) ([]entity.TicketLabel, error) {
	var ticketLabels []entity.TicketLabel
	err := r.db.Where("ticket_id = ?", ticketID).
		Preload("Label").
		Find(&ticketLabels).Error
	return ticketLabels, err
}

func (r *Repository) CheckTicketLabelExists(ticketID, labelID string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.TicketLabel{}).
		Where("ticket_id = ? AND label_id = ?", ticketID, labelID).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) DeleteTicketLabel(ticketID, labelID string) error {
	return r.db.Where("ticket_id = ? AND label_id = ?", ticketID, labelID).
		Delete(&entity.TicketLabel{}).Error
}
