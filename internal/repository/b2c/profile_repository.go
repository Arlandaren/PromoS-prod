package b2c

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"solution/internal/shared/models/b2c"
	"solution/internal/shared/models/b2c/dto"
	"solution/internal/shared/storage/redis"
)

type ProfileRepository interface {
	GetProfile(userID string) (*b2c.User, error)
	UpdateProfile(userID string, req *dto.ProfileUpdateRequest) error
}

type profileRepository struct {
	db  *gorm.DB
	rdb *redis.RDB
}

func NewProfileRepository(db *gorm.DB, rdb *redis.RDB) ProfileRepository {
	return &profileRepository{
		db:  db,
		rdb: rdb,
	}
}

func (r *profileRepository) GetProfile(userID string) (*b2c.User, error) {
	c := context.TODO()

	var user b2c.User
	if err := r.db.WithContext(c).Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *profileRepository) UpdateProfile(userID string, req *dto.ProfileUpdateRequest) error {
	updates := make(map[string]interface{})

	ctx := context.TODO()

	if req.Name != nil {
		updates["name"] = req.Name
	}
	if req.Surname != nil {
		updates["surname"] = req.Surname
	}
	if req.AvatarURL != nil {
		updates["avatar_url"] = req.AvatarURL
	}
	if req.Password != nil {
		updates["password"] = req.Password
	}

	if len(updates) == 0 {
		return dto.ErrNoFieldsToUpdate
	}

	result := r.db.WithContext(ctx).Model(&b2c.User{}).Where("id = ?", userID).Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return dto.ErrNoFieldsToUpdate
	}

	return nil
}
