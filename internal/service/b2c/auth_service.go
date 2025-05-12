package b2c

import (
	"errors"
	"solution/internal/repository/b2c"
	"solution/internal/shared/models/b2c/dto"

	models "solution/internal/shared/models/b2c"
	"solution/internal/shared/utils"
)

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
)

type AuthService interface {
	RegisterUser(req dto.SignUpRequest) (string, string, error)
	AuthenticateUser(req dto.SignInRequest) (string, error)
}

type authService struct {
	repo b2c.AuthRepository
}

func NewAuthService(repo b2c.AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) RegisterUser(req dto.SignUpRequest) (string, string, error) {
	if s.repo.IsEmailRegistered(req.Email) {
		return "", "", ErrEmailAlreadyRegistered
	}

	passwordHash, err := utils.GenerateHashPassword(req.Password)
	if err != nil {
		return "", "", err
	}

	var avatarURL string
	if req.AvatarURL != nil {
		avatarURL = *req.AvatarURL
	}

	newUser := &models.User{
		Name:      req.Name,
		Surname:   req.Surname,
		Email:     req.Email,
		Password:  passwordHash,
		AvatarURL: avatarURL,
		Age:       req.Other.Age,
		Country:   req.Other.Country,
	}

	userID, err := s.repo.CreateUser(newUser)
	if err != nil {
		return "", "", err
	}

	token, err := utils.GenerateToken(userID)
	if err != nil {
		return "", "", err
	}

	err = s.repo.WhitelistToken(token, userID)
	if err != nil {
		return "", "", err
	}

	return token, userID, nil
}

func (s *authService) AuthenticateUser(req dto.SignInRequest) (string, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return "", dto.ErrInvalidCredentials
	}

	if err := utils.CompareHashPassword(req.Password, user.Password); err != nil {
		return "", dto.ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	err = s.repo.WhitelistToken(token, user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
