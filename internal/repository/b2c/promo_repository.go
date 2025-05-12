package b2c

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"solution/internal/shared/models"
	"solution/internal/shared/models/b2b"
	"solution/internal/shared/models/b2c"
	"solution/internal/shared/models/b2c/dto"
	"solution/internal/shared/storage/redis"
	"strings"
	"time"
)

type PromoRepository interface {
	GetPromosForUser(userID string, limit, offset int, category string, active *bool) ([]dto.PromoForUser, int64, error)
	GetPromoByID(id string) (*models.Promo, error)
	GetPromoForUserByID(promoId, userId string) (*dto.PromoForUser, error)
	LikePromo(promoID, userID string) error
	UnlikePromo(promoID, userID string) error
	AddComment(comment *b2c.Comment) error
	GetComments(promoID string, limit, offset int) ([]b2c.Comment, int64, error)
	GetCommentByID(commentID string) (*b2c.Comment, error)
	UpdateComment(comment *b2c.Comment) error
	DeleteComment(commentID string) error
	GetUserByID(userID string) (*b2c.User, error)
}

type promoRepository struct {
	db  *gorm.DB
	rdb *redis.RDB
}

func NewPromoRepository(db *gorm.DB, rdb *redis.RDB) PromoRepository {
	return &promoRepository{
		db:  db,
		rdb: rdb,
	}
}

func (r *promoRepository) GetUserByID(userID string) (*b2c.User, error) {
	var user b2c.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *promoRepository) AddComment(comment *b2c.Comment) error {
	return r.db.Create(comment).Error
}

func (r *promoRepository) GetComments(promoID string, limit, offset int) ([]b2c.Comment, int64, error) {
	var comments []b2c.Comment
	var totalCount int64

	query := r.db.Model(&b2c.Comment{}).Where("promo_id = ?", promoID).Order("created_at DESC")

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, totalCount, nil
}

func (r *promoRepository) GetCommentByID(commentID string) (*b2c.Comment, error) {
	var comment b2c.Comment
	err := r.db.Where("id = ?", commentID).First(&comment).Error
	return &comment, err
}

func (r *promoRepository) UpdateComment(comment *b2c.Comment) error {
	return r.db.Save(comment).Error
}

func (r *promoRepository) DeleteComment(commentID string) error {
	return r.db.Delete(&b2c.Comment{}, "id = ?", commentID).Error
}

func (r *promoRepository) UnlikePromo(promoID, userID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var promo models.Promo
		if err := tx.Where("id = ?", promoID).First(&promo).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.ErrNotFound
			}
			return err
		}

		var userLike b2c.UserLike
		err := tx.Where("user_id = ? AND promo_id = ?", userID, promoID).First(&userLike).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else if err != nil {
			return err
		}

		if err := tx.Delete(&b2c.UserLike{}, "user_id = ? AND promo_id = ?", userID, promoID).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Promo{}).
			Where("id = ? AND like_count > 0", promoID).
			UpdateColumn("like_count", gorm.Expr("like_count - 1")).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r promoRepository) LikePromo(promoID, userID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var userLike b2c.UserLike
		err := tx.Where("user_id = ? AND promo_id = ?", userID, promoID).First(&userLike).Error
		if err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		userLike = b2c.UserLike{
			UserID:  userID,
			PromoID: promoID,
		}
		if err := tx.Create(&userLike).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Promo{}).
			Where("id = ?", promoID).
			UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *promoRepository) GetPromosForUser(userID string, limit, offset int, category string, active *bool) ([]dto.PromoForUser, int64, error) {
	ctx := context.TODO()
	var promos []models.Promo

	// 1. Получаем данные пользователя
	var user b2c.User
	if err := r.db.WithContext(ctx).Model(&b2c.User{}).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, 0, err
	}

	tx := r.db.WithContext(ctx).Model(&models.Promo{})

	// 2. Фильтрация по категории
	if category != "" {
		lowerCategory := strings.ToLower(category)
		tx = tx.Where(`
            EXISTS (
                SELECT 1
                FROM jsonb_array_elements_text(target->'categories') AS c
                WHERE LOWER(c) = ?
            )
        `, lowerCategory)
	}

	// 3. Фильтрация по возрасту пользователя
	if user.Age != 0 {
		tx = tx.Where("COALESCE(target->>'age_from', '0')::int <= ?", user.Age)
		tx = tx.Where("COALESCE(target->>'age_until', '100')::int >= ?", user.Age)
	}

	// 4. Фильтрация по стране пользователя
	if user.Country != "" {
		lowerCountry := strings.ToLower(user.Country)
		tx = tx.Where(`
        (
            target->>'country' IS NULL
            OR LOWER(target->>'country') = ?
        )
    `, lowerCountry)
	}

	if active != nil {
		currentTime := time.Now().UTC()
		if *active {
			// Условия для активных промокодов
			tx = tx.Where(`
            (
                (active_from IS NULL OR active_from <= ?) AND
                (active_until IS NULL OR active_until >= ?) AND
                (
                    (mode = 'COMMON' AND (used_count < max_count OR max_count IS NULL)) OR
                    (mode = 'UNIQUE' AND (used_count < COALESCE(array_length(promo_unique, 1), 0)))
                )
            )
        `, currentTime, currentTime)
		} else {
			// Условия для неактивных промокодов
			tx = tx.Where(`
            NOT (
                (active_from IS NULL OR active_from <= ?) AND
                (active_until IS NULL OR active_until >= ?) AND
                (
                    (mode = 'COMMON' AND (used_count < max_count OR max_count IS NULL)) OR
                    (mode = 'UNIQUE' AND (used_count < COALESCE(array_length(promo_unique, 1), 0)))
                )
            )
        `, currentTime, currentTime)
		}
	}

	// 6. Получаем общее количество промокодов после фильтрации
	var totalCount int64
	if err := tx.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// 7. Применяем сортировку и пагинацию
	tx = tx.Order("created_at DESC").Limit(limit).Offset(offset)

	// 8. Получаем список промокодов
	if err := tx.Find(&promos).Error; err != nil {
		return nil, 0, err
	}

	// Вызываем SetActiveStatus для каждого промокода
	for i := range promos {
		promos[i].SetActiveStatus()
	}

	// 9. Обработка активаций и лайков пользователем
	promoIDs := make([]string, len(promos))
	for i, promo := range promos {
		promoIDs[i] = promo.ID
	}

	// Получение информации об активации промокодов пользователем
	activatedPromoMap, err := r.getActivatedPromos(ctx, userID, promoIDs)
	if err != nil {
		return nil, 0, err
	}

	// Получение информации о лайках промокодов пользователем
	likedPromoMap, err := r.getLikedPromos(ctx, userID, promoIDs)
	if err != nil {
		return nil, 0, err
	}

	// 10. Формирование DTO для ответа
	promoDTOs := make([]dto.PromoForUser, len(promos))
	for i, promo := range promos {
		var company b2b.Company

		result := r.db.WithContext(ctx).Where("id = ?", promo.CompanyID).First(&company)
		if result.Error != nil {
			return nil, 0, result.Error
		}

		var commentCount int64
		result = r.db.Model(&b2c.Comment{}).Where("promo_id = ?", promo.ID).Count(&commentCount)
		if result.Error != nil {
			return nil, 0, result.Error
		}

		promoDTOs[i] = dto.PromoForUser{
			PromoID:           promo.ID,
			CompanyID:         promo.CompanyID,
			CompanyName:       company.Name,
			Description:       promo.Description,
			ImageURL:          promo.ImageURL,
			Active:            promo.Active,
			IsActivatedByUser: activatedPromoMap[promo.ID],
			LikeCount:         promo.LikeCount,
			IsLikedByUser:     likedPromoMap[promo.ID],
			CommentCount:      commentCount,
		}
	}

	return promoDTOs, totalCount, nil
}

func (r *promoRepository) GetPromoForUserByID(promoId, userId string) (*dto.PromoForUser, error) {
	ctx := context.Background()
	var promo models.Promo
	err := r.db.WithContext(ctx).Where("id = ?", promoId).First(&promo).Error
	if err != nil {
		return nil, err
	}

	promo.SetActiveStatus()

	promoIDs := make([]string, 0)

	promoIDs = append(promoIDs, promo.ID)

	activatedPromoMap, err := r.getActivatedPromos(ctx, userId, promoIDs)
	if err != nil {
		return nil, err
	}

	//3 Получение информации о лайках промокодов пользователем
	likedPromoMap, err := r.getLikedPromos(ctx, userId, promoIDs)
	if err != nil {
		return nil, err
	}

	// 4. Формирование DTO для ответа

	var company b2b.Company

	result := r.db.WithContext(ctx).Where("id = ?", promo.CompanyID).First(&company)
	if result.Error != nil {
		return nil, result.Error
	}

	var commentCount int64
	result = r.db.Model(&b2c.Comment{}).Where("promo_id = ?", promo.ID).Count(&commentCount)
	if result.Error != nil {
		return nil, result.Error
	}

	promoDTO := dto.PromoForUser{
		PromoID:           promo.ID,
		CompanyID:         promo.CompanyID,
		CompanyName:       company.Name,
		Description:       promo.Description,
		ImageURL:          promo.ImageURL,
		Active:            promo.Active,
		IsActivatedByUser: activatedPromoMap[promo.ID],
		LikeCount:         promo.LikeCount,
		IsLikedByUser:     likedPromoMap[promo.ID],
		CommentCount:      commentCount,
	}

	return &promoDTO, nil
}

func (r *promoRepository) GetPromoByID(id string) (*models.Promo, error) {
	ctx := context.TODO()
	var promo models.Promo
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&promo).Error
	if err != nil {
		return nil, err
	}

	return &promo, nil
}

func (r *promoRepository) getActivatedPromos(ctx context.Context, userID string, promoIDs []string) (map[string]bool, error) {
	var activatedPromoIDs []string
	err := r.db.WithContext(ctx).
		Model(&models.PromoActivation{}).
		Where("user_id = ? AND promo_id IN ?", userID, promoIDs).
		Pluck("promo_id", &activatedPromoIDs).Error
	if err != nil {
		return nil, err
	}
	activatedPromoMap := make(map[string]bool)
	for _, id := range activatedPromoIDs {
		activatedPromoMap[id] = true
	}
	return activatedPromoMap, nil
}

func (r *promoRepository) getLikedPromos(ctx context.Context, userID string, promoIDs []string) (map[string]bool, error) {
	var likedPromoIDs []string
	err := r.db.WithContext(ctx).
		Model(&b2c.UserLike{}).
		Where("user_id = ? AND promo_id IN ?", userID, promoIDs).
		Pluck("promo_id", &likedPromoIDs).Error
	if err != nil {
		return nil, err
	}
	likedPromoMap := make(map[string]bool)
	for _, id := range likedPromoIDs {
		likedPromoMap[id] = true
	}
	return likedPromoMap, nil
}
