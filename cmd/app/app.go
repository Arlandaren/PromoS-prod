package app

import (
	"context"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"

	b2b_repo "solution/internal/repository/b2b"

	b2b_service "solution/internal/service/b2b"

	b2c_repo "solution/internal/repository/b2c"

	b2c_service "solution/internal/service/b2c"

	di "solution/internal/service/services"
	"solution/internal/shared/config"
	"solution/internal/shared/storage/postgres"
	"solution/internal/shared/storage/redis"
	server "solution/internal/transport/http"
	"sync"
	"syscall"
)

type App struct {
	cfg         *config.Config
	db          *gorm.DB
	redisClient *redis.RDB
	appServer   *server.Server
	wg          *sync.WaitGroup
}

func NewApp() (*App, error) {

	cfg, err := config.Init()
	if err != nil {
		return nil, err
	}

	db, err := postgres.InitPostgres(cfg.Postgres)
	if err != nil {
		return nil, err
	}

	redisClient, err := redis.NewRedisClient(cfg)
	if err != nil {
		return nil, err
	}

	if err := registerEnv(cfg, redisClient); err != nil {
		return nil, err
	}

	if err := registerRepositories(db, redisClient); err != nil {
		return nil, err
	}

	if err := registerServices(cfg); err != nil {
		return nil, err
	}

	var appServer *server.Server
	if err := di.GetService(&appServer); err != nil {
		return nil, err
	}

	return &App{
		cfg:         cfg,
		db:          db,
		redisClient: redisClient,
		appServer:   appServer,
		wg:          &sync.WaitGroup{},
	}, nil
}

func registerEnv(cfg *config.Config, redisClient *redis.RDB) error {
	envRegistrations := map[string]func() error{
		"config": func() error { return di.AddSingleton(func() *config.Config { return cfg }) },
		"redis":  func() error { return di.AddSingleton(func() *redis.RDB { return redisClient }) },
	}
	for name, register := range envRegistrations {
		if err := register(); err != nil {
			log.Fatalf("Failed to register %s: %v", name, err)
			return err
		}
		log.Println(name, "registered")
	}
	return nil
}

func registerRepositories(db *gorm.DB, redisClient *redis.RDB) error {
	repoRegistrations := map[string]func() error{
		"b2bAuthRepo": func() error {
			return di.AddSingleton(func() b2b_repo.AuthRepository { return b2b_repo.NewAuthRepository(db, redisClient) })
		},
		"b2bPromoRepo": func() error {
			return di.AddSingleton(func() b2b_repo.PromoRepository { return b2b_repo.NewPromoRepository(db, redisClient) })
		},
		"b2cAuthRepo": func() error {
			return di.AddSingleton(func() b2c_repo.AuthRepository { return b2c_repo.NewAuthRepository(db, redisClient) })
		},
		"b2cProfileRepo": func() error {
			return di.AddSingleton(func() b2c_repo.ProfileRepository { return b2c_repo.NewProfileRepository(db, redisClient) })
		},
		"b2cPromoRepo": func() error {
			return di.AddSingleton(func() b2c_repo.PromoRepository { return b2c_repo.NewPromoRepository(db, redisClient) })
		},
	}
	for name, register := range repoRegistrations {
		if err := register(); err != nil {
			log.Fatalf("Failed to register %s: %v", name, err)
			return err
		}
		log.Println(name, "registered")
	}
	return nil
}

func registerServices(cfg *config.Config) error {
	serviceRegistrations := map[string]func() error{
		"b2bAuthService": func() error {
			return di.AddSingleton(func(repo b2b_repo.AuthRepository) b2b_service.AuthService { return b2b_service.NewAuthService(repo) })
		},
		"b2bPromoService": func() error {
			return di.AddSingleton(func(repo b2b_repo.PromoRepository) b2b_service.PromoService { return b2b_service.NewPromoService(repo) })
		},
		"b2cAuthService": func() error {
			return di.AddSingleton(func(repo b2c_repo.AuthRepository) b2c_service.AuthService { return b2c_service.NewAuthService(repo) })
		},
		"b2cProfileService": func() error {
			return di.AddSingleton(func(repo b2c_repo.ProfileRepository) b2c_service.ProfileService {
				return b2c_service.NewProfileService(repo)
			})
		},
		"b2cPromoService": func() error {
			return di.AddSingleton(func(repo b2c_repo.PromoRepository) b2c_service.PromoService { return b2c_service.NewPromoService(repo) })
		},
	}
	for name, register := range serviceRegistrations {
		if err := register(); err != nil {
			log.Fatalf("Failed to register %s: %v", name, err)
			return err
		}
		log.Println(name, "registered")
	}

	err := di.AddSingleton(func(cfg *config.Config) *server.Server {
		return server.NewServer(cfg)
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := a.appServer.StartHttpServer(ctx); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down HTTP server...")

	cancel()
	a.wg.Wait()
	log.Println("Redis gracefully stopped.")
}
