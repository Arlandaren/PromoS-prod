package models

import "time"

type PromoActivation struct {
	ID          string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PromoID     string `gorm:"type:uuid"`
	UserID      string `gorm:"type:uuid"`
	ActivatedAt time.Time
}
