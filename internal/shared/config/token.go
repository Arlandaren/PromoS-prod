package config

import (
	"os"
)

func GetJwtKey() []byte {
	jwtKey := os.Getenv("RANDOM_SECRET")
	if jwtKey == "" {
		panic("RANDOM_SECRET is not set in the environment")
	}

	return []byte(jwtKey)
}
