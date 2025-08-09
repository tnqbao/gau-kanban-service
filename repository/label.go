package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
)

func (r *Repository) CreateLabel(label *entity.Label) error {
	return r.db.Create(label).Error
}

func (r *Repository) GetAllLabel() ([]entity.Label, error) {
	var labels []entity.Label
	err := r.db.Order("name ASC").Find(&labels).Error
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

func (r *Repository) GetLabelByTicketID(ticketID string) ([]entity.Label, error) {
	var labels []entity.Label
	err := r.db.Table("labels").
		Joins("JOIN ticket_labels ON labels.id = ticket_labels.label_id").
		Where("ticket_labels.ticket_id = ?", ticketID).
		Find(&labels).Error
	return labels, err
}
