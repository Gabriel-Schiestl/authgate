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

func NewUserInfo(props UserInfoProps) (UserInfo, *exceptions.BusinessException) {
	if props.UserID == "" {
		return nil, exceptions.NewBusinessException("user ID cannot be empty")
	}

	return &userInfo{
		userID: props.UserID,
		name:   props.Name,
		roles:  props.Roles,
	}, nil
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