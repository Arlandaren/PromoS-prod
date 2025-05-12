package b2b

import (
	"errors"
	"gorm.io/gorm"
	b2b2 "solution/internal/repository/b2b"
	"solution/internal/shared/models"
	"solution/internal/shared/models/b2b/dto"
)

type PromoService interface {
	CreatePromo(req dto.PromoCreateRequest) (string, error)
	GetPromos(companyID string, limit, offset int, sortBy string, country []string) ([]models.Promo, int64, error)
	GetPromoByID(companyID string, promoID string) (*dto.PromoReadOnlyResponse, error)
	UpdatePromo(companyID string, promoID string, req dto.PromoPatchRequest) (*models.Promo, error)
	GetPromoStatByID(companyID string, promoID string) (*dto.PromoStatResponse, error)
}

type promoService struct {
	repo b2b2.PromoRepository
}

func NewPromoService(repo b2b2.PromoRepository) PromoService {
	return &promoService{repo: repo}
}

func (s *promoService) CreatePromo(req dto.PromoCreateRequest) (string, error) {
	return s.repo.CreatePromo(req)
}

func (s *promoService) GetPromos(companyID string, limit, offset int, sortBy string, country []string) ([]models.Promo, int64, error) {
	return s.repo.GetPromos(companyID, limit, offset, sortBy, country)
}

func (s *promoService) GetPromoByID(companyID string, promoID string) (*dto.PromoReadOnlyResponse, error) {
	promo, err := s.repo.GetPromoByID(promoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	company, err := s.repo.GetCompanyById(promo.CompanyID)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, gorm.ErrRecordNotFound
	}

	promoResponse := &dto.PromoReadOnlyResponse{
		Description: promo.Description,
		ImageURL:    promo.ImageURL,
		Target:      promo.Target,
		MaxCount:    promo.MaxCount,
		ActiveFrom:  promo.ActiveFrom,
		ActiveUntil: promo.ActiveUntil,
		Mode:        promo.Mode,
		PromoCommon: promo.PromoCommon,
		PromoUnique: promo.PromoUnique,
		PromoId:     promo.ID,
		CompanyID:   promo.CompanyID,
		CompanyName: company.Name,
		LikeCount:   promo.LikeCount,
		UsedCount:   promo.UsedCount,
		Active:      promo.Active,
	}

	if promo.CompanyID != companyID {
		return nil, dto.ErrorNoAccess
	}

	return promoResponse, nil
}

func (s *promoService) UpdatePromo(companyID string, promoID string, req dto.PromoPatchRequest) (*models.Promo, error) {
	promo, err := s.repo.GetPromoByID(promoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	if req.MaxCount != nil {
		if promo.Mode == "UNIQUE" && *req.MaxCount > 1 {
			return nil, dto.ErrBadRequest
		}
	}

	if promo.CompanyID != companyID {
		return nil, dto.ErrorNoAccess
	}

	return s.repo.UpdatePromo(promoID, req)
}

func (s *promoService) GetPromoStatByID(companyID string, promoID string) (*dto.PromoStatResponse, error) {
	// First, check if the promo exists and belongs to the company
	promo, err := s.repo.GetPromoByID(promoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrorPromoNotFound
		}
		return nil, err
	}

	if promo.CompanyID != companyID {
		return nil, dto.ErrorNoAccessToPromo
	}

	promoStat, err := s.repo.GetPromoStatByID(promoID)
	if err != nil {
		return nil, err
	}

	return promoStat, nil
}
