package b2b

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"log"

	redis "solution/internal/shared/storage/redis"

	redisPkg "github.com/go-redis/redis/v8"
	"solution/internal/shared/models/b2b"
	"solution/internal/shared/models/b2b/dto"
	"time"
)

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidCredentials     = errors.New("invalid credentials")
)

type AuthRepository interface {
	CreateCompany(req dto.SignUpRequest) (string, error)
	GetCompany(email string) (*b2b.Company, error)
	IsEmailRegistered(email string) bool
	WhitelistToken(token, id string) error
	ValidateToken(userId string) (string, error)
}

type authRepository struct {
	db  *gorm.DB
	rdb *redis.RDB
}

func NewAuthRepository(db *gorm.DB, rdb *redis.RDB) AuthRepository {
	return &authRepository{
		db:  db,
		rdb: rdb,
	}
}

func (r *authRepository) CreateCompany(req dto.SignUpRequest) (string, error) {
	ctx := context.TODO()

	if r.IsEmailRegistered(req.Email) {
		return "", ErrEmailAlreadyRegistered
	}

	company := b2b.Company{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := r.db.WithContext(ctx).Create(&company).Error; err != nil {
		return "", err
	}

	return company.ID, nil
}

func (r *authRepository) GetCompany(email string) (*b2b.Company, error) {
	ctx := context.TODO()

	var company b2b.Company
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&company).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	return &company, nil
}

func (r *authRepository) IsEmailRegistered(email string) bool {
	ctx := context.TODO()

	var exists bool
	err := r.db.WithContext(ctx).Model(&b2b.Company{}).Select("count(*) > 0").Where("email = ?", email).Scan(&exists).Error
	if err != nil {
		return false
	}

	return exists
}

func (r *authRepository) WhitelistToken(token, id string) error {
	err := r.rdb.Client.Set(context.TODO(), id, token, 24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *authRepository) ValidateToken(userId string) (string, error) {
	ctx := context.TODO()
	token, err := r.rdb.Client.Get(ctx, userId).Result()
	if errors.Is(err, redisPkg.Nil) {
		log.Println("Token does not exist for the given userId")
		return "", nil
	} else if err != nil {
		log.Println("Error while fetching token from Redis:", err)
		return "", err
	}
	return token, nil
}
