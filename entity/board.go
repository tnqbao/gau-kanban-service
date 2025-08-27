package entity

type Board struct {
	ID          string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title       string `gorm:"type:text;not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Archived    bool   `gorm:"type:boolean;default:false" json:"archived"`
	CreatedAt   string `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt   string `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`
}

func (Board) TableName() string {
	return "boards"
}
