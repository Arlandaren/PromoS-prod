package b2b

import (
	"errors"
	"solution/internal/repository/b2b"
	"solution/internal/shared/models/b2b/dto"
	"solution/internal/shared/utils"
)

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidCredentials     = errors.New("invalid credentials")
)

type AuthService interface {
	RegisterCompany(req dto.SignUpRequest) (string, string, error)
	AuthenticateCompany(req dto.SignInRequest) (string, error)
}

type authService struct {
	repo b2b.AuthRepository
}

func NewAuthService(repo b2b.AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) RegisterCompany(req dto.SignUpRequest) (string, string, error) {
	if s.repo.IsEmailRegistered(req.Email) {
		return "", "", ErrEmailAlreadyRegistered
	}

	hash, err := utils.GenerateHashPassword(req.Password)
	if err != nil {
		return "", "", err
	}
	req.Password = hash

	companyID, err := s.repo.CreateCompany(req)
	if err != nil {
		return "", "", err
	}

	token, err := utils.GenerateToken(companyID)
	if err != nil {
		return "", "", err
	}

	err = s.repo.WhitelistToken(token, companyID)
	if err != nil {
		return "", "", err
	}

	return token, companyID, nil
}

func (s *authService) AuthenticateCompany(req dto.SignInRequest) (string, error) {
	company, err := s.repo.GetCompany(req.Email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := utils.CompareHashPassword(req.Password, company.Password); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(company.ID)
	if err != nil {
		return "", err
	}

	err = s.repo.WhitelistToken(token, company.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
