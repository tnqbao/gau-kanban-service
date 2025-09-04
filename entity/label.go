package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Label struct {
	ID          string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BoardID     string    `json:"board_id" gorm:"type:uuid;not null"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null"`
	Color       string    `json:"color" gorm:"type:varchar(7);default:'#6B7280'"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (l *Label) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}

type TicketLabel struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TicketID  string    `json:"ticket_id" gorm:"type:uuid;not null"`
	LabelID   string    `json:"label_id" gorm:"type:uuid;not null"`
	Label     Label     `json:"label" gorm:"foreignKey:LabelID"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (tl *TicketLabel) BeforeCreate(tx *gorm.DB) error {
	if tl.ID == "" {
		tl.ID = uuid.New().String()
	}
	return nil
}
