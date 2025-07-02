package dtos

import "github.com/Gabriel-Schiestl/authgate/internal/src/domain/models"

type LoginDTO struct {
	IdentifierType  models.IdentifierType `json:"identifier_type"`
	IdentifierValue string              `json:"identifier_value"`
	Password        string              `json:"password"`
}

type LoginResponseDTO struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	UserInfo     UserInfoDTO `json:"user_info"`
}