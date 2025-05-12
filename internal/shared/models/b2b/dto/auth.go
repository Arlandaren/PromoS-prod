package dto

import (
	"errors"
	"regexp"
	"unicode/utf8"
)

var (
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrInvalidPassword    = errors.New("password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	ErrInvalidCompanyName = errors.New("company name must be between 5 and 50 characters long")
)

type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	CompanyID string `json:"company_id"`
}

func (req *SignUpRequest) Validate() error {
	if utf8.RuneCountInString(req.Name) < 5 || utf8.RuneCountInString(req.Name) > 50 {
		return ErrInvalidCompanyName
	}

	if err := validateEmail(req.Email); err != nil {
		return err
	}

	if err := validatePassword(req.Password); err != nil {
		return err
	}

	return nil
}

func (req *SignInRequest) Validate() error {
	if err := validateEmail(req.Email); err != nil {
		return err
	}

	if err := validatePassword(req.Password); err != nil {
		return err
	}

	return nil
}

func validateEmail(email string) error {
	if utf8.RuneCountInString(email) < 8 || utf8.RuneCountInString(email) > 120 {
		return ErrInvalidEmail
	}
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !re.MatchString(email) {
		return ErrInvalidEmailFormat
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
