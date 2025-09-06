package entity

import (
	"encoding/json"
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

// Custom JSON structure for assignee response
type TicketAssigneeResponse struct {
	ID       string `json:"id"`
	TicketID string `json:"ticket_id"`
	MemberID string `json:"member_id"`
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
}

func (ta *TicketAssignee) BeforeCreate(tx *gorm.DB) error {
	if ta.ID == "" {
		ta.ID = uuid.New().String()
	}
	return nil
}

// MarshalJSON customizes the JSON output to only include user_id and full_name
func (ta TicketAssignee) MarshalJSON() ([]byte, error) {
	response := TicketAssigneeResponse{
		ID:       ta.ID,
		TicketID: ta.TicketID,
		MemberID: ta.MemberID,
		UserID:   ta.Member.UserID,
		FullName: ta.Member.User.FullName,
	}
	return json.Marshal(response)
}
