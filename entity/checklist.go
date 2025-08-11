package entity

type Checklist struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TicketID  string `gorm:"type:uuid;not null" json:"ticket_id"`
	Title     string `gorm:"type:text;not null" json:"title"`
	Completed bool   `gorm:"type:boolean;default:false" json:"completed"`
	Position  int    `gorm:"type:integer;default:0" json:"position"`
	CreatedAt string `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt string `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`
}

func (Checklist) TableName() string {
	return "checklists"
}
