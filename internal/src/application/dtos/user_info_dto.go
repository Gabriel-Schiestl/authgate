package dtos

type UserInfoDTO struct {
	UserID string   `json:"user_id"`
	Name   string   `json:"name"`
	Roles  []string `json:"roles"`
}