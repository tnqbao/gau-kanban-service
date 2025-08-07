package entity

//CREATE TABLE IF NOT EXISTS labels (
//id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//name TEXT NOT NULL,
//color TEXT NOT NULL,
//created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
//updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
//);

type Label struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name      string `gorm:"type:text;not null" json:"name"`
	Color     string `gorm:"type:text;not null" json:"color"`
	CreatedAt string `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt string `gorm:"type:timestamp with time zone;default:now()" json:"updated_at"`
}

func (Label) TableName() string {
	return "labels"
}
