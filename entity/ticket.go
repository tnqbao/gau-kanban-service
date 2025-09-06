package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ticket struct {
	ID           string     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TicketNumber string     `json:"ticket_number" gorm:"type:varchar(10);unique;not null"`
	Title        string     `json:"title" gorm:"type:varchar(255);not null"`
	Description  string     `json:"description" gorm:"type:text"`
	ColumnID     string     `json:"column_id" gorm:"type:uuid;not null"`
	Position     int        `json:"position" gorm:"default:0"`
	Priority     string     `json:"priority" gorm:"type:varchar(20);default:'medium'"`
	DueDate      *time.Time `json:"due_date" gorm:"type:timestamp"`
	Status       string     `json:"status" gorm:"type:varchar(50);default:'open'"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Assignees  []TicketAssignee `json:"assignees" gorm:"foreignKey:TicketID"`
	Labels     []TicketLabel    `json:"labels" gorm:"foreignKey:TicketID"`
	Checklists []Checklist      `json:"checklists" gorm:"foreignKey:TicketID"`
}

func (t *Ticket) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

// MarshalJSON ensures that nil slices are marshaled as empty arrays
func (t Ticket) MarshalJSON() ([]byte, error) {
	type Alias Ticket

	// Initialize nil slices as empty arrays
	if t.Assignees == nil {
		t.Assignees = []TicketAssignee{}
	}
	if t.Labels == nil {
		t.Labels = []TicketLabel{}
	}
	if t.Checklists == nil {
		t.Checklists = []Checklist{}
	}

	return json.Marshal((Alias)(t))
}
