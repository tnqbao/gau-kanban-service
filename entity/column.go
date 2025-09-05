package entity

import (
	"encoding/json"
	"time"
)

type Column struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	BoardID   string    `gorm:"type:uuid;not null" json:"board_id"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Position  int       `gorm:"type:integer;not null" json:"position"`
	WipLimit  *int      `gorm:"type:integer" json:"wip_limit"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`

	// Relationships
	Tickets []Ticket `gorm:"foreignKey:ColumnID" json:"tickets"`
}

func (Column) TableName() string {
	return "columns"
}

// MarshalJSON ensures that nil slices are marshaled as empty arrays
func (c Column) MarshalJSON() ([]byte, error) {
	type Alias Column

	// Initialize nil slices as empty arrays
	if c.Tickets == nil {
		c.Tickets = []Ticket{}
	}

	return json.Marshal((Alias)(c))
}
