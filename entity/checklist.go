package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Checklist struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TicketID  string    `json:"ticket_id" gorm:"type:uuid;not null"`
	Title     string    `json:"title" gorm:"type:varchar(255);not null"`
	Status    string    `json:"status" gorm:"type:varchar(20);default:'pending'"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (cl *Checklist) BeforeCreate(tx *gorm.DB) error {
	if cl.ID == "" {
		cl.ID = uuid.New().String()
	}
	return nil
}
