package entity

type TicketComment struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TicketID  string `gorm:"type:uuid;not null" json:"ticket_id"`
	UserID    string `gorm:"type:uuid;not null" json:"user_id"`
	Content   string `gorm:"type:text;not null" json:"content"`
	CreatedAt string `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
}

func (TicketComment) TableName() string {
	return "ticket_comments"
}
