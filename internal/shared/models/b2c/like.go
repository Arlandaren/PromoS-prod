package b2c

type UserLike struct {
	UserID  string `gorm:"type:uuid;primaryKey"`
	PromoID string `gorm:"type:uuid;primaryKey"`
}
