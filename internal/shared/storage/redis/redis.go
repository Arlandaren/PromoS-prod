package redis

import (
	"context"
	"solution/internal/shared/config"

	"github.com/go-redis/redis/v8"
)

type RDB struct {
	Client *redis.Client
}

func NewRedisClient(config *config.Config) (*RDB, error) {
	address := config.Redis.ConnStr

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})

	if err := client.Ping(context.TODO()).Err(); err != nil {
		return &RDB{}, err
	}

	return &RDB{
		Client: client,
	}, nil
}
