package entity

import "time"

type Member struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"type:uuid;not null" json:"user_id"`
	BoardID   string    `gorm:"type:uuid;not null" json:"board_id"`
	Role      string    `gorm:"type:varchar(50);default:'member'" json:"role"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`
}

func (Member) TableName() string {
	return "members"
}
