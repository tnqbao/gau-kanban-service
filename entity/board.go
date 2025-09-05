package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Board struct {
	ID          string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text"`
	OwnerID     string    `json:"owner_id" gorm:"type:uuid;not null"`
	Archived    bool      `json:"archived" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations - these are loaded separately and not stored in DB
	Columns []Column `json:"columns" gorm:"-"`
	Members []Member `json:"members" gorm:"-"`
	Labels  []Label  `json:"labels" gorm:"-"`
}

func (b *Board) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

// MarshalJSON ensures that nil slices are marshaled as empty arrays
func (b Board) MarshalJSON() ([]byte, error) {
	type Alias Board

	// Initialize nil slices as empty arrays
	if b.Columns == nil {
		b.Columns = []Column{}
	}
	if b.Members == nil {
		b.Members = []Member{}
	}
	if b.Labels == nil {
		b.Labels = []Label{}
	}

	return json.Marshal((Alias)(b))
}
