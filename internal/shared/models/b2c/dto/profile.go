package dto

import (
	"unicode/utf8"
)

type ProfileResponse struct {
	Name      string             `json:"name" binding:"required"`
	Surname   string             `json:"surname" binding:"required"`
	Email     string             `json:"email" binding:"required"`
	AvatarURL string             `json:"avatar_url"`
	Other     UserTargetSettings `json:"other" binding:"required"`
}

type ProfileUpdateRequest struct {
	Name      *string `json:"name,omitempty"`
	Surname   *string `json:"surname,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Password  *string `json:"password,omitempty"`
}

func (req *ProfileUpdateRequest) Validate() error {
	if req.Name != nil {
		if utf8.RuneCountInString(*req.Name) < 1 || utf8.RuneCountInString(*req.Name) > 100 {
			return ErrInvalidName
		}
	}

	if req.Surname != nil {
		if utf8.RuneCountInString(*req.Surname) < 1 || utf8.RuneCountInString(*req.Surname) > 120 {
			return ErrInvalidSurname
		}
	}

	if req.Password != nil {
		if err := validatePassword(*req.Password); err != nil {
			return err
		}
	}

	if req.AvatarURL != nil {
		if err := validateAvatarURL(*req.AvatarURL); err != nil {
			return err
		}
	}

	return nil
}
