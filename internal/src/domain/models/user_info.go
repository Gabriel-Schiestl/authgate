package models

import "github.com/Gabriel-Schiestl/go-clarch/v2/domain/exceptions"

type UserInfo interface {
	GetUserID() string
	GetName() string
	GetRoles() []string
}

type userInfo struct {
	userID string
	name   string
	roles  []string
}

type UserInfoProps struct {
	UserID string
	Name   string
	Roles  []string
}

func NewUserInfo(props UserInfoProps) (*exceptions.BusinessException, UserInfo) {
	if props.UserID == "" {
		return exceptions.NewBusinessException("user ID cannot be empty"), nil
	}

	return nil, &userInfo{
		userID: props.UserID,
		name:   props.Name,
		roles:  props.Roles,
	}
}

func (u *userInfo) GetUserID() string {
	return u.userID
}
func (u *userInfo) GetName() string {
	return u.name
}
func (u *userInfo) GetRoles() []string {
	return u.roles
}