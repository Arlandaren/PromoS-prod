package config

import (
	"errors"
	"os"
)

type Postgres struct {
	ConnStr  string
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

func getPostgres() (*Postgres, error) {
	connStr := os.Getenv("POSTGRES_CONN")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	database := os.Getenv("POSTGRES_DATABASE")
	username := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")

	if connStr == "" || host == "" || port == "" || database == "" || username == "" || password == "" {
		return nil, errors.New("one or more environment variables are missing")
	}

	return &Postgres{
		ConnStr:  connStr,
		Host:     host,
		Port:     port,
		Database: database,
		Username: username,
		Password: password,
	}, nil
}
