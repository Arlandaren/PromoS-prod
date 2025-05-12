package b2c

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	repo "solution/internal/repository/b2c"
	"solution/internal/shared/models/b2c/dto"
)

type ProfileService interface {
	GetProfile(userId string) (*dto.ProfileResponse, error)
	UpdateProfile(userId string, req dto.ProfileUpdateRequest) error
}

type profileService struct {
	repo repo.ProfileRepository
}

func NewProfileService(repo repo.ProfileRepository) ProfileService {
	return &profileService{repo: repo}
}

func (s *profileService) GetProfile(userId string) (*dto.ProfileResponse, error) {
	user, err := s.repo.GetProfile(userId)
	if err != nil {
		return nil, err
	}

	fmt.Println(user)

	return &dto.ProfileResponse{
		Name:      user.Name,
		Surname:   user.Surname,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		Other:     dto.UserTargetSettings{Age: user.Age, Country: user.Country},
	}, nil
}

func (s *profileService) UpdateProfile(userId string, req dto.ProfileUpdateRequest) error {
	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		newPassword := string(hashedPassword)
		req.Password = &newPassword
	}
	return s.repo.UpdateProfile(userId, &req)
}
