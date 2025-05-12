package b2b

import (
	"context"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"solution/internal/shared/models"
	"solution/internal/shared/models/b2b"
	"solution/internal/shared/models/b2b/dto"
	"solution/internal/shared/storage/redis"
	"sort"
	"strings"
)

type PromoRepository interface {
	CreatePromo(req dto.PromoCreateRequest) (string, error)
	GetPromos(companyID string, limit, offset int, sortBy string, country []string) ([]models.Promo, int64, error)
	GetPromoByID(promoID string) (*models.Promo, error)
	UpdatePromo(promoID string, req dto.PromoPatchRequest) (*models.Promo, error)
	GetPromoStatByID(promoID string) (*dto.PromoStatResponse, error)
	GetCompanyById(id string) (*b2b.Company, error)
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

func (r *promoRepository) CreatePromo(req dto.PromoCreateRequest) (string, error) {
	ctx := context.TODO()

	promo := models.Promo{
		CompanyID:   req.CompanyID,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Mode:        req.Mode,
		PromoCommon: req.PromoCommon,
		PromoUnique: pq.StringArray(req.PromoUnique),
		Target:      *req.Target,
		MaxCount:    *req.MaxCount,
		ActiveFrom:  req.ActiveFrom,
		ActiveUntil: req.ActiveUntil,
	}

	if err := r.db.WithContext(ctx).Create(&promo).Error; err != nil {
		log.Printf("Error creating promo: %v", err)
		return "", err
	}

	return promo.ID, nil
}

func (r *promoRepository) GetPromos(companyID string, limit, offset int, sortBy string, country []string) ([]models.Promo, int64, error) {
	ctx := context.TODO()
	var promos []models.Promo

	fmt.Println(companyID)

	tx := r.db.Debug().WithContext(ctx).Model(&models.Promo{}).Where("company_id = ?", companyID)

	if len(country) > 0 {
		lowerCountries := make([]string, len(country))
		for i, c := range country {
			lowerCountries[i] = strings.ToLower(c)
		}

		tx = tx.Where(
			"(LOWER(target->>'country') IN ? OR target->>'country' IS NULL OR target->>'country' = '')",
			lowerCountries,
		)
	}

	if sortBy != "" {
		if sortBy == "active_from" || sortBy == "active_until" {
			tx = tx.Order(clause.OrderByColumn{
				Column: clause.Column{Name: sortBy},
				Desc:   true,
			})
		}
	} else {
		tx = tx.Order("created_at DESC")
	}

	var totalCount int64
	if err := tx.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	tx = tx.Limit(limit).Offset(offset)

	if err := tx.Find(&promos).Error; err != nil {
		return nil, 0, err
	}

	fmt.Println(promos)

	return promos, totalCount, nil
}

func (r *promoRepository) GetPromoByID(promoID string) (*models.Promo, error) {
	ctx := context.TODO()
	var promo models.Promo

	if err := r.db.WithContext(ctx).First(&promo, "id = ?", promoID).Error; err != nil {
		return nil, err
	}

	return &promo, nil
}

func (r *promoRepository) UpdatePromo(promoID string, req dto.PromoPatchRequest) (*models.Promo, error) {
	ctx := context.TODO()
	var promo models.Promo

	if err := r.db.WithContext(ctx).First(&promo, "id = ?", promoID).Error; err != nil {
		return nil, err
	}

	if req.Description != "" {
		promo.Description = req.Description
	}
	if req.ImageURL != "" {
		promo.ImageURL = req.ImageURL
	}
	if req.Target != nil {
		promo.Target = *req.Target
	}
	if req.MaxCount != nil {
		promo.MaxCount = *req.MaxCount
	}
	if req.ActiveFrom != nil {
		promo.ActiveFrom = req.ActiveFrom
	}
	if req.ActiveUntil != nil {
		promo.ActiveUntil = req.ActiveUntil
	}

	if err := r.db.WithContext(ctx).Save(&promo).Error; err != nil {
		log.Printf("Error updating promo: %v", err)
		return nil, err
	}

	return &promo, nil
}

func (r *promoRepository) GetPromoStatByID(promoID string) (*dto.PromoStatResponse, error) {
	ctx := context.TODO()
	var promo models.Promo
	var countryActivations []dto.CountryActivation

	// Получаем промокод по ID из таблицы Promo, включая поле used_count
	if err := r.db.WithContext(ctx).
		Select("used_count").
		First(&promo, "id = ?", promoID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrorPromoNotFound
		}
		return nil, err
	}

	// Получаем количество активаций по странам из таблицы PromoActivation
	if err := r.db.WithContext(ctx).
		Model(&models.PromoActivation{}).
		Select("users.country as country_code, COUNT(*) as activations_count").
		Joins("JOIN users ON promo_activations.user_id = users.id").
		Where("promo_activations.promo_id = ?", promoID).
		Group("users.country").
		Scan(&countryActivations).Error; err != nil {
		return nil, err
	}

	// Проверяем, есть ли активации по странам
	if len(countryActivations) == 0 {
		countryActivations = []dto.CountryActivation{}
	}

	// Сортируем страны по коду (регистронезависимо)
	sort.Slice(countryActivations, func(i, j int) bool {
		return strings.ToLower(countryActivations[i].Country.Code) < strings.ToLower(countryActivations[j].Country.Code)
	})

	promoStat := &dto.PromoStatResponse{
		ActivationsCount: promo.UsedCount,
		Countries:        countryActivations,
	}

	return promoStat, nil
}

func (r *promoRepository) GetCompanyById(id string) (*b2b.Company, error) {
	ctx := context.TODO()
	var company b2b.Company

	if err := r.db.WithContext(ctx).First(&company, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrNotFound
		}
		return nil, err
	}

	return &company, nil
}
