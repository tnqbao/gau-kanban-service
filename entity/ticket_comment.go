package entity

//id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
//user_id UUID NOT NULL,
//content TEXT NOT NULL,
//created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
//);

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
