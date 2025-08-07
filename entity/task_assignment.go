package entity

//CREATE TABLE IF NOT EXISTS task_assignments (
//id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
//user_id UUID NOT NULL,
//assigned_at TIMESTAMP WITH TIME ZONE DEFAULT now()
//);

type TaskAssignment struct {
	ID           string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TicketID     string `gorm:"type:uuid;not null" json:"ticket_id"`
	UserID       string `gorm:"type:uuid;not null" json:"user_id"`
	UserFullName string `gorm:"type:text;not null" json:"full_name"`
	AssignedAt   string `gorm:"type:timestamp with time zone;default:now()" json:"assigned_at"`
}

func (TaskAssignment) TableName() string {
	return "task_assignments"
}
