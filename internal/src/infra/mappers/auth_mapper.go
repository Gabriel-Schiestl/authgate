package mappers

import (
	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/models"
	"github.com/Gabriel-Schiestl/authgate/internal/src/infra/entities"
)

func ModelToDomain(entity entities.Auth) (models.Auth, error) {
	userInfo, err := userInfoModelToDomain(entity.UserInfo)
	if err != nil {
		return nil, err
	}

	domain, domainErr := models.LoadAuth(models.AuthProps{
		ID:                 entity.ID,
		IdentifierType:     IdentifierTypeToDomain(entity.IdentifierType),
		IdentifierValue:    entity.IdentifierValue,
		Password:           entity.Password,
		UserInfo:           userInfo,
		EncryptToken:       entity.EncryptToken,
		LastLoginAt:        entity.LastLoginAt,
		WrongAttempts:      entity.WrongAttempts,
		MaxWrongAttempts:   &entity.MaxWrongAttempts,
		RecoveryToken:      entity.RecoveryToken,
		MaxTokenAgeSeconds: &entity.MaxTokenAgeSeconds,
	})
	if domainErr != nil {
		return nil, domainErr
	}

	return domain, nil
}

func DomainToModel(domain models.Auth) entities.Auth {
	return entities.Auth{
		ID:                 domain.GetID(),
		IdentifierType:     IdentifierTypeFromDomain(domain.GetIdentifierType()),
		IdentifierValue:    domain.GetIdentifierValue(),
		Password:           domain.GetPassword(),
		UserInfo: userInfoDomainToModel(domain.GetUserInfo()),
		EncryptToken:       domain.GetEncryptToken(),
		LastLoginAt:        domain.GetLastLoginAt(),
		WrongAttempts:      domain.GetWrongAttempts(),
		MaxWrongAttempts:   *domain.GetMaxWrongAttempts(),
		RecoveryToken:      domain.GetRecoveryToken(),
		MaxTokenAgeSeconds: *domain.GetMaxTokenAgeSeconds(),
	}
}

func userInfoModelToDomain(entity entities.UserInfo) (models.UserInfo, error) {
	model, err := models.NewUserInfo(
		models.UserInfoProps{
			UserID: entity.UserID,
			Name:   entity.Name,
			Roles:  entity.Roles,
		},
	)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func userInfoDomainToModel(domain models.UserInfo) entities.UserInfo {
	return entities.UserInfo{
		UserID: domain.GetUserID(),
		Name:   domain.GetName(),
		Roles:  domain.GetRoles(),
	}
}