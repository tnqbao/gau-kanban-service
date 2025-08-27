package entity

type Column struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	BoardID   string `gorm:"type:uuid;not null" json:"board_id"`
	Title     string `gorm:"type:text;not null" json:"title"`
	Position  int    `gorm:"type:integer;autoIncrement" json:"position"`
	WipLimit  *int   `gorm:"type:integer" json:"wip_limit"`
	CreatedAt string `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt string `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`
}

func (Column) TableName() string {
	return "columns"
}
