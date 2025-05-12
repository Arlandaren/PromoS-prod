package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"solution/internal/shared/models/b2b"
	"strings"
	"time"
	"unicode/utf8"
)

type Promo struct {
	ID          string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id,omitempty"`
	CompanyID   string         `json:"company_id,omitempty"`
	Description string         `json:"description"`
	ImageURL    string         `json:"image_url,omitempty"`
	Mode        string         `json:"mode"`
	PromoCommon string         `json:"promo_common,omitempty"`
	PromoUnique pq.StringArray `json:"promo_unique,omitempty" gorm:"type:text[]"`
	Target      Target         `json:"target" gorm:"type:jsonb"`
	MaxCount    int            `json:"max_count"`
	ActiveFrom  *b2b.Date      `json:"active_from,omitempty"`
	ActiveUntil *b2b.Date      `json:"active_until,omitempty"`
	LikeCount   int            `json:"like_count,required"`
	UsedCount   int            `json:"used_count,required"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"-"`
	Active      bool           `gorm:"-" json:"active"`
}

func (p *Promo) SetActiveStatus() {
	currentTime := time.Now().UTC()
	isActive := true

	if p.ActiveFrom != nil && currentTime.Before(p.ActiveFrom.Time) {
		isActive = false
	}

	if p.ActiveUntil != nil && currentTime.After(p.ActiveUntil.Time) {
		isActive = false
	}

	if p.Mode == "COMMON" {
		if p.MaxCount == 0 || p.UsedCount >= p.MaxCount {
			isActive = false
		}
	} else if p.Mode == "UNIQUE" {
		promoUniqueCount := len(p.PromoUnique)
		if promoUniqueCount == 0 || p.UsedCount >= promoUniqueCount {
			isActive = false
		}
	}

	p.Active = isActive
}

type Target struct {
	AgeFrom    *int     `json:"age_from,omitempty"`
	AgeUntil   *int     `json:"age_until,omitempty"`
	Country    string   `json:"country,omitempty"`
	Categories []string `json:"categories,omitempty"`
}

func (t *Target) Validate() error {
	// Validate AgeFrom
	if t.AgeFrom != nil {
		if *t.AgeFrom < 0 || *t.AgeFrom > 100 {
			return errors.New("age_from must be between 0 and 100")
		}
	}

	// Validate AgeUntil
	if t.AgeUntil != nil {
		if *t.AgeUntil < 0 || *t.AgeUntil > 100 {
			return errors.New("age_until must be between 0 and 100")
		}
	}

	// If both AgeFrom and AgeUntil are provided, validate their relationship
	if t.AgeFrom != nil && t.AgeUntil != nil {
		if *t.AgeFrom > *t.AgeUntil {
			return errors.New("age_from cannot be greater than age_until")
		}
	}

	// Validate Country
	if t.Country != "" {
		countryCode := strings.ToUpper(t.Country)
		if _, valid := ValidCountryCodes[countryCode]; !valid {
			return fmt.Errorf("country must be a valid ISO 3166-1 alpha-2 code, got '%s'", t.Country)
		}
	}

	// Validate Categories
	if t.Categories != nil && len(t.Categories) > 0 {
		for _, category := range t.Categories {
			if utf8.RuneCountInString(category) < 2 || utf8.RuneCountInString(category) > 20 {
				return errors.New("each category must be between 2 and 20 characters long")
			}
		}
	}

	return nil
}

func (t Target) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *Target) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, t)
}
