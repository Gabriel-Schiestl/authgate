package dtos

import "github.com/Gabriel-Schiestl/authgate/internal/src/domain/models"

type RegisterDTO struct {
	IdentifierType     models.IdentifierType `json:"identifier_type"`
	IdentifierValue    string                `json:"identifier_value"`
	Password           string                `json:"password"`
	UserInfo           UserInfoDTO           `json:"user_info"`
	EncryptToken       bool                  `json:"encrypt_token"`
	MaxTokenAgeSeconds *int                  `json:"max_token_age_seconds,omitempty"`
	MaxWrongAttempts   *int                  `json:"max_wrong_attempts,omitempty"`
}

type RegisterResponseDTO struct {
	IdentifierType  models.IdentifierType      `json:"identifier_type"`
	IdentifierValue string      `json:"identifier_value"`
	UserInfo        UserInfoDTO `json:"user_info"`
}
