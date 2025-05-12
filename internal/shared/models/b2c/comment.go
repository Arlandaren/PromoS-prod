package b2c

import (
	"time"
)

type Comment struct {
	ID        string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    string `gorm:"type:uuid;not null"`
	PromoID   string `gorm:"type:uuid;not null"`
	Text      string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	User User `gorm:"foreignKey:UserID"`
}
