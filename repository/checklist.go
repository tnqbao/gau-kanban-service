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
	err := r.db.Where("ticket_id = ?", ticketID).Order("order ASC").Find(&checklists).Error
	return checklists, err
}

func (r *Repository) GetNextChecklistOrder(ticketID string) (int, error) {
	var maxOrder int
	err := r.db.Model(&entity.Checklist{}).
		Where("ticket_id = ?", ticketID).
		Select("COALESCE(MAX(order), 0) + 1").
		Scan(&maxOrder).Error
	return maxOrder, err
}

func (r *Repository) UpdateChecklist(checklist *entity.Checklist) error {
	return r.db.Save(checklist).Error
}

func (r *Repository) DeleteChecklist(id string) error {
	return r.db.Delete(&entity.Checklist{}, "id = ?", id).Error
}

func (r *Repository) GetChecklistByID(id string) (*entity.Checklist, error) {
	var checklist entity.Checklist
	err := r.db.Where("id = ?", id).First(&checklist).Error
	if err != nil {
		return nil, err
	}
	return &checklist, nil
}
