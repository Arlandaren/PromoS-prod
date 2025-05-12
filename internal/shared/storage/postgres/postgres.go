package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"solution/internal/shared/config"
	"solution/internal/shared/models"
	"solution/internal/shared/models/b2b"
	"solution/internal/shared/models/b2c"
)

func InitPostgres(cfg *config.Postgres) (*gorm.DB, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return nil, err
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	err = db.AutoMigrate(&b2c.User{}, &b2b.Company{}, &models.Promo{}, &models.PromoActivation{}, &b2c.UserLike{}, &b2c.Comment{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
		return nil, err
	}
	return db, nil
}
