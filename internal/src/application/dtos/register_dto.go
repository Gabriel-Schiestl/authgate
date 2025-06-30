package dtos

type RegisterDTO struct {
	IdentifierType     int32       `json:"identifier_type"`
	IdentifierValue    string      `json:"identifier_value"`
	Password           string      `json:"password"`
	UserInfo           UserInfoDTO `json:"user_info"`
	EncryptToken       bool        `json:"encrypt_token"`
	MaxTokenAgeSeconds *int        `json:"max_token_age_seconds,omitempty"`
	MaxWrongAttempts   *int        `json:"max_wrong_attempts,omitempty"`
}

type RegisterResponseDTO struct {
	IdentifierType  string      `json:"identifier_type"`
	IdentifierValue string      `json:"identifier_value"`
	UserInfo        UserInfoDTO `json:"user_info"`
}
