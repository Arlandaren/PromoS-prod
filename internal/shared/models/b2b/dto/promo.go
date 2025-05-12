package dto

import (
	"errors"
	"net/url"
	models2 "solution/internal/shared/models"
	models "solution/internal/shared/models/b2b"
	"unicode/utf8"
)

var (
	ErrInvalidPromoMode     = errors.New("promo mode must be either 'COMMON' or 'UNIQUE'")
	ErrInvalidPromoCode     = errors.New("promo code must be between 5 and 30 characters long")
	ErrInvalidPromoMaxCount = errors.New("max count must be greater than 0")
	ErrPromoUniqueTooShort  = errors.New("each promo_unique must be between 3 and 30 characters long")
	ErrPromoUniqueRequired  = errors.New("at least one promo_unique is required for UNIQUE mode")
	ErrPromoCommonRequired  = errors.New("promo_common is required for COMMON mode")
	ErrPromoUniqueNotUnique = errors.New("promo_unique values must be unique")
	ErrorPromoNotFound      = errors.New("promo not found")
	ErrorNoAccess           = errors.New("no access to this resource")
	ErrorNoAccessToPromo    = errors.New("no access to promo")
)

type Country struct {
	Code string `json:"code"`
}

type CountryActivation struct {
	Country     Country `json:"country"`
	Activations int     `json:"activations_count"`
}

type PromoStatResponse struct {
	ActivationsCount int                 `json:"activations_count"`
	Countries        []CountryActivation `json:"countries,omitempty"`
}

type PromoCreateRequest struct {
	CompanyID   string          `json:"company_id"`
	Description string          `json:"description" binding:"required"`
	ImageURL    string          `json:"image_url"`
	Mode        string          `json:"mode" binding:"required"`
	PromoCommon string          `json:"promo_common,omitempty"`
	PromoUnique []string        `json:"promo_unique,omitempty"`
	Target      *models2.Target `json:"target" binding:"required"`
	MaxCount    *int            `json:"max_count" binding:"required"`
	ActiveFrom  *models.Date    `json:"active_from,omitempty"`
	ActiveUntil *models.Date    `json:"active_until,omitempty"`
}

type PromoReadOnlyResponse struct {
	Description string         `json:"description" binding:"required"`
	ImageURL    string         `json:"image_url"`
	Target      models2.Target `json:"target" binding:"required"` // Use models.Target
	MaxCount    int            `json:"max_count" binding:"required"`
	ActiveFrom  *models.Date   `json:"active_from"`
	ActiveUntil *models.Date   `json:"active_until"`
	Mode        string         `json:"mode" binding:"required"`
	PromoCommon string         `json:"promo_common,omitempty"`
	PromoUnique []string       `json:"promo_unique,omitempty"`
	PromoId     string         `json:"promo_id" binding:"required"`
	CompanyID   string         `json:"company_id" binding:"required"`
	CompanyName string         `json:"company_name" binding:"required"`
	LikeCount   int            `json:"like_count" binding:"required"`
	UsedCount   int            `json:"used_count" binding:"required"`
	Active      bool           `json:"active" binding:"required"`
}

type PromoPatchRequest struct {
	Description string          `json:"description,omitempty" `
	ImageURL    string          `json:"image_url,omitempty"`
	Target      *models2.Target `json:"target,omitempty"`
	MaxCount    *int            `json:"max_count,omitempty"`
	ActiveFrom  *models.Date    `json:"active_from,omitempty"`
	ActiveUntil *models.Date    `json:"active_until,omitempty"`
}

func (req *PromoPatchRequest) Validate() error {
	if req.Description != "" {
		length := utf8.RuneCountInString(req.Description)
		if length < 10 || length > 300 {
			return errors.New("description must be between 10 and 300 characters long")
		}
	}

	if req.ImageURL != "" {
		if len(req.ImageURL) > 350 {
			return errors.New("image_url must be at most 350 characters long")
		}
		if _, err := url.ParseRequestURI(req.ImageURL); err != nil {
			return errors.New("image_url must be a valid URL")
		}
	}

	if req.Target != nil {
		if err := req.Target.Validate(); err != nil {
			return err
		}
	}

	if req.MaxCount != nil {
		if *req.MaxCount < 0 || *req.MaxCount > 100000000 {
			return errors.New("max_count must be between 0 and 100,000,000")
		}
	}

	if req.ActiveFrom != nil && req.ActiveUntil != nil {
		if req.ActiveFrom.Time.After(req.ActiveUntil.Time) {
			return errors.New("active_from cannot be after active_until")
		}
	}

	return nil
}

func (req *PromoCreateRequest) Validate() error {
	if err := validatePromoMode(req.Mode); err != nil {
		return err
	}

	if req.Description != "" {
		length := utf8.RuneCountInString(req.Description)
		if length < 10 || length > 300 {
			return errors.New("description must be between 10 and 300 characters long")
		}
	}
	if req.ImageURL != "" {
		if len(req.ImageURL) > 350 {
			return errors.New("image_url must be at most 350 characters long")
		}
		if _, err := url.ParseRequestURI(req.ImageURL); err != nil {
			return errors.New("image_url must be a valid URL")
		}
	}

	switch req.Mode {
	case "COMMON":
		if req.PromoCommon == "" {
			return ErrPromoCommonRequired
		}

		if utf8.RuneCountInString(req.PromoCommon) < 5 || utf8.RuneCountInString(req.PromoCommon) > 30 {
			return ErrInvalidPromoCode
		}

	case "UNIQUE":
		if len(req.PromoUnique) == 0 {
			return ErrPromoUniqueRequired
		}

		for _, code := range req.PromoUnique {
			if utf8.RuneCountInString(code) < 3 || utf8.RuneCountInString(code) > 30 {
				return ErrPromoUniqueTooShort
			}
		}

		uniqueCodes := make(map[string]struct{})
		for _, code := range req.PromoUnique {
			if _, exists := uniqueCodes[code]; exists {
				return ErrPromoUniqueNotUnique
			}
			uniqueCodes[code] = struct{}{}
		}

	default:
		return ErrInvalidPromoMode
	}
	if req.MaxCount != nil {
		if req.Mode == "UNIQUE" && *req.MaxCount > 1 {
			return ErrBadRequest
		}
	}

	if req.MaxCount != nil {
		if *req.MaxCount < 0 || *req.MaxCount > 100000000 {
			return errors.New("max_count must be between 0 and 100,000,000")
		}
	}

	if req.ActiveFrom != nil && req.ActiveUntil != nil {
		if req.ActiveFrom.Time.After(req.ActiveUntil.Time) {
			return errors.New("active_from cannot be after active_until")
		}
	}

	if req.Target != nil {
		if err := req.Target.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func validatePromoMode(mode string) error {
	if mode != "COMMON" && mode != "UNIQUE" {
		return ErrInvalidPromoMode
	}
	return nil
}
