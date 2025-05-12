package b2b

type Company struct {
	ID       string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name     string `gorm:"size:100;not null"`
	Email    string `gorm:"size:100;not null;unique"`
	Password string `gorm:"size:100;not null"`
}
