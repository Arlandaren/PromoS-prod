package config

import (
	"errors"
	"os"
)

type Redis struct {
	ConnStr string
}

func getRedis() (*Redis, error) {
	redisConn := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	if redisConn == ":" { // оба значения пустые
		return nil, errors.New("not found REDIS_HOST or REDIS_PORT")
	}
	return &Redis{ConnStr: redisConn}, nil
}
