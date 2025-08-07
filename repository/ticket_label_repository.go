package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
	"gorm.io/gorm"
)

type TicketLabelRepository struct {
	db *gorm.DB
}

func NewTicketLabelRepository(db *gorm.DB) TicketLabelRepositoryInterface {
	return &TicketLabelRepository{db: db}
}

func (r *TicketLabelRepository) AddLabelToTicket(ticketID, labelID string) error {
	ticketLabel := entity.TicketLabel{
		TicketID: ticketID,
		LabelID:  labelID,
	}
	return r.db.Create(&ticketLabel).Error
}

func (r *TicketLabelRepository) RemoveLabelFromTicket(ticketID, labelID string) error {
	return r.db.Where("ticket_id = ? AND label_id = ?", ticketID, labelID).Delete(&entity.TicketLabel{}).Error
}

func (r *TicketLabelRepository) GetTicketsByLabelID(labelID string) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	err := r.db.Table("tickets").
		Joins("JOIN ticket_labels ON tickets.id = ticket_labels.ticket_id").
		Where("ticket_labels.label_id = ?", labelID).
		Find(&tickets).Error
	return tickets, err
}

func (r *TicketLabelRepository) GetLabelsByTicketID(ticketID string) ([]entity.Label, error) {
	var labels []entity.Label
	err := r.db.Table("labels").
		Joins("JOIN ticket_labels ON labels.id = ticket_labels.label_id").
		Where("ticket_labels.ticket_id = ?", ticketID).
		Find(&labels).Error
	return labels, err
}

func (r *TicketLabelRepository) RemoveAllLabelsFromTicket(ticketID string) error {
	return r.db.Where("ticket_id = ?", ticketID).Delete(&entity.TicketLabel{}).Error
}
