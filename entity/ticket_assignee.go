package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketAssignee struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TicketID  string    `json:"ticket_id" gorm:"type:uuid;not null"`
	MemberID  string    `json:"member_id" gorm:"type:uuid;not null"`
	Member    Member    `json:"member" gorm:"foreignKey:MemberID"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (ta *TicketAssignee) BeforeCreate(tx *gorm.DB) error {
	if ta.ID == "" {
		ta.ID = uuid.New().String()
	}
	return nil
}
