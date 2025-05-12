package config

type Config struct {
	Postgres *Postgres
	Redis    *Redis
	Server   *Server
}

func Init() (*Config, error) {
	postgresConfig, err := getPostgres()
	if err != nil {
		return nil, err
	}

	redisCfg, err := getRedis()
	if err != nil {
		return nil, err
	}

	serverCfg, err := getServer()
	if err != nil {
		return nil, err
	}

	return &Config{
		Postgres: postgresConfig,
		Redis:    redisCfg,
		Server:   serverCfg,
	}, nil
}
