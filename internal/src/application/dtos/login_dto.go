package dtos

type LoginDTO struct {
	IdentifierType  int    `json:"identifier_type"`
	IdentifierValue string `json:"identifier_value"`
	Password        string `json:"password"`
}

type LoginResponseDTO struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	UserInfo     UserInfoDTO `json:"user_info"`
}