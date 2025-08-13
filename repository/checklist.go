package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
)

// Checklist methods

func (r *Repository) CreateChecklist(checklist *entity.Checklist) error {
	return r.db.Create(checklist).Error
}

func (r *Repository) GetChecklistsByTicketID(ticketID string) ([]entity.Checklist, error) {
	var checklists []entity.Checklist
	err := r.db.Where("ticket_id = ?", ticketID).Order("position ASC").Find(&checklists).Error
	return checklists, err
}

func (r *Repository) GetChecklistByID(id string) (*entity.Checklist, error) {
	var checklist entity.Checklist
	err := r.db.Where("id = ?", id).First(&checklist).Error
	if err != nil {
		return nil, err
	}
	return &checklist, nil
}

func (r *Repository) UpdateChecklist(checklist *entity.Checklist) error {
	return r.db.Save(checklist).Error
}

func (r *Repository) DeleteChecklist(id string) error {
	return r.db.Where("id = ?", id).Delete(&entity.Checklist{}).Error
}

func (r *Repository) DeleteChecklistsByTicketID(ticketID string) error {
	return r.db.Where("ticket_id = ?", ticketID).Delete(&entity.Checklist{}).Error
}

func (r *Repository) UpdateChecklistPosition(id string, position int) error {
	return r.db.Model(&entity.Checklist{}).Where("id = ?", id).Update("position", position).Error
}

func (r *Repository) GetMaxChecklistPosition(ticketID string) (int, error) {
	var maxPosition int
	err := r.db.Model(&entity.Checklist{}).Where("ticket_id = ?", ticketID).Select("COALESCE(MAX(position), 0)").Scan(&maxPosition).Error
	return maxPosition, err
}
