package models

import (
	"time"

	"github.com/Gabriel-Schiestl/go-clarch/v2/domain/exceptions"
	"github.com/Gabriel-Schiestl/go-clarch/v2/utils"
)

type Auth interface {
	GetID() string
	GetIdentifierType() IdentifierType
	GetIdentifierValue() string
	GetPassword() string
	GetUserInfo() UserInfo
	GetEncryptToken() bool
	GetLastLoginAt() *time.Time
	GetWrongAttempts() int
	GetMaxWrongAttempts() *int
	GetRecoveryToken() *string
	GetMaxTokenAgeSeconds() *int
}

type auth struct {
	id    string
	identifierType IdentifierType
	identifierValue string
	password string
	userInfo UserInfo
	encryptToken bool
	lastLoginAt *time.Time
	wrongAttempts int
	maxWrongAttempts *int
	recoveryToken *string
	maxTokenAgeSeconds *int
}

type AuthProps struct {
	ID    string
	IdentifierType IdentifierType
	IdentifierValue string
	Password string
	UserInfo UserInfo
	EncryptToken bool
	LastLoginAt *time.Time
	WrongAttempts int
	MaxWrongAttempts *int
	RecoveryToken *string
	MaxTokenAgeSeconds *int
}

func NewAuth(props AuthProps) (Auth, *exceptions.BusinessException) {
	if props.IdentifierType == 0 {
		return nil, exceptions.NewBusinessException("identifier type cannot be unspecified")
	}
	if props.IdentifierValue == "" {
		return nil, exceptions.NewBusinessException("identifier value cannot be empty")
	}
	if props.Password == "" {
		return nil, exceptions.NewBusinessException("password cannot be empty")
	}
	if props.UserInfo == nil {
		return nil, exceptions.NewBusinessException("user info cannot be nil")
	}

	newAuth := &auth{
		id:              props.ID,
		identifierType:  props.IdentifierType,
		identifierValue: props.IdentifierValue,
		password:        props.Password,
		userInfo:        props.UserInfo,
		encryptToken:    props.EncryptToken,
		lastLoginAt:     props.LastLoginAt,
		wrongAttempts:   props.WrongAttempts,
		maxWrongAttempts: props.MaxWrongAttempts,
		recoveryToken:   props.RecoveryToken,
		maxTokenAgeSeconds: props.MaxTokenAgeSeconds,
	}
	
	if newAuth.id == "" {
		newAuth.id = utils.GenerateUUID()
	}
	if newAuth.maxWrongAttempts == nil {
		defaultMaxWrongAttempts := 5
		newAuth.maxWrongAttempts = &defaultMaxWrongAttempts
	}
	if newAuth.maxTokenAgeSeconds == nil {
		defaultMaxTokenAgeSeconds := 604800
		newAuth.maxTokenAgeSeconds = &defaultMaxTokenAgeSeconds
	}

	return newAuth, nil
}

func LoadAuth(props AuthProps) (Auth, *exceptions.BusinessException) {
	return NewAuth(props)
}

func (a *auth) GetID() string {
	return a.id
}

func (a *auth) GetIdentifierType() IdentifierType {
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

func (a *auth) GetRecoveryToken() *string {
	return a.recoveryToken
}

func (a *auth) GetMaxTokenAgeSeconds() *int {
	return a.maxTokenAgeSeconds
}