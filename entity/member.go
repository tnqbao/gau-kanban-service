package entity

type Member struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	BoardID   string `gorm:"type:uuid;not null" json:"board_id"`
	MemberID  string `gorm:"type:uuid;not null" json:"member_id"`
	FullName  string `gorm:"type:text;not null" json:"full_name"`
	CreatedAt string `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt string `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`
}

func (Member) TableName() string {
	return "members"
}
