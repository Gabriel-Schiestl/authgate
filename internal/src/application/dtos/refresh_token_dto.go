package dtos

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponseDTO struct {
	AccessToken string      `json:"access_token"`
	UserInfo    UserInfoDTO `json:"user_info"`
}