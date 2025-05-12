package config

import (
	"errors"
	"os"
	"strings"
)

type Server struct {
	Addr string
	Port string
}

func getServer() (*Server, error) {
	serverConf := os.Getenv("SERVER_ADDRESS")

	if serverConf == "" {
		return nil, errors.New("not found REDIS_HOST or REDIS_PORT")
	}
	port := strings.Split(serverConf, ":")[1]

	return &Server{Addr: serverConf,
		Port: port,
	}, nil
}
