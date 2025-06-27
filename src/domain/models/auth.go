package models

import (
	"time"

	"github.com/Gabriel-Schiestl/go-clarch/v2/domain/exceptions"
	"github.com/Gabriel-Schiestl/go-clarch/v2/utils"
)

type Auth interface {
	GetID() string
	GetIdentifierType() string
	GetIdentifierValue() string
	GetPassword() string
	GetUserInfo() UserInfo
	GetEncryptToken() bool
	GetLastLoginAt() *time.Time
	GetWrongAttempts() int
	GetMaxWrongAttempts() *int
	GetRecoveryToken() string
	GetMaxTokenAgeSeconds() *int
	GetCreatedAt() *time.Time
}

type auth struct {
	id    string
	identifierType string
	identifierValue string
	password string
	userInfo UserInfo
	encryptToken bool
	lastLoginAt *time.Time
	wrongAttempts int
	maxWrongAttempts *int
	recoveryToken string
	maxTokenAgeSeconds *int
	createdAt *time.Time
}

type AuthProps struct {
	id    string
	identifierType IdentifierType
	identifierValue string
	password string
	userInfo UserInfo
	encryptToken bool
	lastLoginAt *time.Time
	wrongAttempts int
	maxWrongAttempts *int
	recoveryToken string
	maxTokenAgeSeconds *int
	createdAt *time.Time
}

func NewAuth(props AuthProps) (*exceptions.BusinessException, Auth) {
	if props.identifierType == 0 {
		return exceptions.NewBusinessException("identifier type cannot be unspecified"), nil
	}
	if props.identifierValue == "" {
		return exceptions.NewBusinessException("identifier value cannot be empty"), nil
	}
	if props.password == "" {
		return exceptions.NewBusinessException("password cannot be empty"), nil
	}
	if props.userInfo == nil {
		return exceptions.NewBusinessException("user info cannot be nil"), nil
	}

	identifier := IdentifierType_name[int32(props.identifierType)]

	newAuth := &auth{
		id:              props.id,
		identifierType:  identifier,
		identifierValue: props.identifierValue,
		password:        props.password,
		userInfo:        props.userInfo,
		encryptToken:    props.encryptToken,
		lastLoginAt:     props.lastLoginAt,
		wrongAttempts:   props.wrongAttempts,
		maxWrongAttempts: props.maxWrongAttempts,
		recoveryToken:   props.recoveryToken,
		createdAt:       props.createdAt,
	}

	if newAuth.id == "" {
		newAuth.id = utils.GenerateUUID()
	}
	if newAuth.createdAt == nil {
		now := time.Now()
		newAuth.createdAt = &now
	}
	if newAuth.maxWrongAttempts == nil {
		defaultMaxWrongAttempts := 5
		newAuth.maxWrongAttempts = &defaultMaxWrongAttempts
	}
	if newAuth.maxTokenAgeSeconds == nil {
		defaultMaxTokenAgeSeconds := 604800
		newAuth.maxTokenAgeSeconds = &defaultMaxTokenAgeSeconds
	}

	return nil, newAuth
}

func (a *auth) GetID() string {
	return a.id
}

func (a *auth) GetIdentifierType() string {
	return a.identifierType
}

func (a *auth) GetIdentifierValue() string {
	return a.identifierValue
}

func (a *auth) GetPassword() string {
	return a.password
}

func (a *auth) GetUserInfo() UserInfo {
	return a.userInfo
}

func (a *auth) GetEncryptToken() bool {
	return a.encryptToken
}

func (a *auth) GetLastLoginAt() *time.Time {
	return a.lastLoginAt
}

func (a *auth) GetWrongAttempts() int {
	return a.wrongAttempts
}

func (a *auth) GetMaxWrongAttempts() *int {
	return a.maxWrongAttempts
}

func (a *auth) GetRecoveryToken() string {
	return a.recoveryToken
}

func (a *auth) GetMaxTokenAgeSeconds() *int {
	return a.maxTokenAgeSeconds
}

func (a *auth) GetCreatedAt() *time.Time {
	return a.createdAt
}