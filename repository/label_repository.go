package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
	"gorm.io/gorm"
)

type LabelRepository struct {
	db *gorm.DB
}

func NewLabelRepository(db *gorm.DB) LabelRepositoryInterface {
	return &LabelRepository{db: db}
}

func (r *LabelRepository) Create(label *entity.Label) error {
	return r.db.Create(label).Error
}

func (r *LabelRepository) GetAll() ([]entity.Label, error) {
	var labels []entity.Label
	err := r.db.Order("name ASC").Find(&labels).Error
	return labels, err
}

func (r *LabelRepository) GetByID(id string) (*entity.Label, error) {
	var label entity.Label
	err := r.db.Where("id = ?", id).First(&label).Error
	if err != nil {
		return nil, err
	}
	return &label, nil
}

func (r *LabelRepository) Update(label *entity.Label) error {
	return r.db.Save(label).Error
}

func (r *LabelRepository) Delete(id string) error {
	return r.db.Delete(&entity.Label{}, "id = ?", id).Error
}

func (r *LabelRepository) GetByTicketID(ticketID string) ([]entity.Label, error) {
	var labels []entity.Label
	err := r.db.Table("labels").
		Joins("JOIN ticket_labels ON labels.id = ticket_labels.label_id").
		Where("ticket_labels.ticket_id = ?", ticketID).
		Find(&labels).Error
	return labels, err
}
