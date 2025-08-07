package entity

type Column struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name      string `gorm:"type:text;not null" json:"name"`
	Position  int    `gorm:"type:integer;not null;default:0" json:"position"`
	CreatedAt string `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt string `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`
}

func (Column) TableName() string {
	return "columns"
}
