package entity

type Ticket struct {
	ID          string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TicketNo    string `gorm:"type:text;unique;not null" json:"ticket_no"`
	ColumnID    string `gorm:"type:uuid;not null" json:"column_id"`
	Title       string `gorm:"type:text;not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	DueDate     string `gorm:"type:date" json:"due_date"`
	Priority    string `gorm:"type:text" json:"priority"`
	Position    int    `gorm:"type:integer;default:0" json:"position"`
	CreatedAt   string `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt   string `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`
}

func (Ticket) TableName() string {
	return "tickets"
}
