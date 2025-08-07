package entity

type TicketLabel struct {
	TicketID string `gorm:"type:uuid;not null" json:"ticket_id"`
	LabelID  string `gorm:"type:uuid;not null" json:"label_id"`
}
