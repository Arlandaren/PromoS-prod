package b2c

import (
	"context"
	"errors"
	redisPkg "github.com/go-redis/redis/v8"
	"log"
	models "solution/internal/shared/models/b2c"
	"solution/internal/shared/storage/redis"
	"time"

	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
)

type AuthRepository interface {
	CreateUser(user *models.User) (string, error)
	GetUserByEmail(email string) (*models.User, error)
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

func (r *authRepository) CreateUser(user *models.User) (string, error) {
	ctx := context.TODO()

	if r.IsEmailRegistered(user.Email) {
		return "", ErrEmailAlreadyRegistered
	}

	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return "", err
	}

	return user.ID, nil
}

func (r *authRepository) GetUserByEmail(email string) (*models.User, error) {
	ctx := context.TODO()

	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *authRepository) IsEmailRegistered(email string) bool {
	ctx := context.TODO()

	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false
	}

	return count > 0
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
