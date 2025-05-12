package dto

import (
	"errors"
	"fmt"
	"regexp"
	"solution/internal/shared/models"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalidName        = errors.New("name must be between 1 and 100 characters long")
	ErrInvalidSurname     = errors.New("surname must be between 1 and 120 characters long")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidPassword    = errors.New("password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	ErrInvalidAvatarURL   = errors.New("avatar_url must be a valid URL and up to 350 characters long")
	ErrInvalidAge         = errors.New("age must be between 0 and 100")
	ErrInvalidCountry     = errors.New("country must be provided")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type SignUpRequest struct {
	Name      string             `json:"name" binding:"required"`
	Surname   string             `json:"surname" binding:"required"`
	Email     string             `json:"email" binding:"required"`
	Password  string             `json:"password" binding:"required"`
	AvatarURL *string            `json:"avatar_url"`
	Other     UserTargetSettings `json:"other" binding:"required"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

func (req *SignUpRequest) Validate() error {
	if utf8.RuneCountInString(req.Name) < 1 || utf8.RuneCountInString(req.Name) > 100 {
		return ErrInvalidName
	}

	if utf8.RuneCountInString(req.Surname) < 1 || utf8.RuneCountInString(req.Surname) > 120 {
		return ErrInvalidSurname
	}

	if err := validateEmail(req.Email); err != nil {
		return err
	}

	if err := validatePassword(req.Password); err != nil {
		return err
	}

	// Validate AvatarURL if it was provided
	if req.AvatarURL != nil {
		if err := validateAvatarURL(*req.AvatarURL); err != nil {
			return err
		}
	}

	if req.Other == (UserTargetSettings{}) {
		return errors.New("other is required")
	}

	if err := req.Other.Validate(); err != nil {
		return err
	}

	return nil
}

type UserTargetSettings struct {
	Age     int    `json:"age" binding:"required"`
	Country string `json:"country" binding:"required"`
}

func (settings *UserTargetSettings) Validate() error {
	if settings.Age < 0 || settings.Age > 100 {
		return ErrInvalidAge
	}

	if settings.Country != "" {
		countryCode := strings.ToUpper(settings.Country)
		if _, valid := models.ValidCountryCodes[countryCode]; !valid {
			return fmt.Errorf("country must be a valid ISO 3166-1 alpha-2 code, got '%s'", settings.Country)
		}
	}

	return nil
}

func (req *SignInRequest) Validate() error {
	if err := validateEmail(req.Email); err != nil {
		return err
	}

	if err := validatePassword(req.Password); err != nil {
		return ErrInvalidPassword
	}

	if req.Password == "" {
		return ErrInvalidPassword
	}

	return nil
}

func validateEmail(email string) error {
	if utf8.RuneCountInString(email) < 8 || utf8.RuneCountInString(email) > 120 {
		return ErrInvalidEmail
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !re.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

func validatePassword(password string) error {
	if utf8.RuneCountInString(password) < 8 || utf8.RuneCountInString(password) > 60 {
		return ErrInvalidPassword
	}

	upper := regexp.MustCompile(`[A-Z]`)
	if !upper.MatchString(password) {
		return ErrInvalidPassword
	}

	lower := regexp.MustCompile(`[a-z]`)
	if !lower.MatchString(password) {
		return ErrInvalidPassword
	}

	number := regexp.MustCompile(`\d`)
	if !number.MatchString(password) {
		return ErrInvalidPassword
	}

	special := regexp.MustCompile(`[@$!%*?&]`)
	if !special.MatchString(password) {
		return ErrInvalidPassword
	}

	return nil
}

func validateAvatarURL(url string) error {
	if utf8.RuneCountInString(url) == 0 || utf8.RuneCountInString(url) > 350 {
		return ErrInvalidAvatarURL
	}

	re := regexp.MustCompile(`^(http|https):\/\/[^\s]+$`)
	if !re.MatchString(url) {
		return ErrInvalidAvatarURL
	}

	return nil
}
