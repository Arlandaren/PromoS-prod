package b2c

type User struct {
	ID        string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string `gorm:"size:100;not null"`
	Surname   string `gorm:"size:120;not null"`
	Email     string `gorm:"size:100;not null;unique"`
	Password  string `gorm:"size:255;not null"`
	AvatarURL string `gorm:"size:350"`
	Age       int    `gorm:"not null"`
	Country   string `gorm:"size:100;not null"`
}
