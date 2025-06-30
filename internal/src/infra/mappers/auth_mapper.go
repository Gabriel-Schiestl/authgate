package mappers

import (
	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/models"
	"github.com/Gabriel-Schiestl/authgate/internal/src/infra/entities"
)

func ModelToDomain(entity entities.Auth) (models.Auth, error) {
	domain, err := models.LoadAuth(models.AuthProps{
		ID:                 entity.ID,
		IdentifierType:     IdentifierTypeToDomain(entity.IdentifierType),
		IdentifierValue:    entity.IdentifierValue,
		Password:           entity.Password,
		EncryptToken:       entity.EncryptToken,
		LastLoginAt:        entity.LastLoginAt,
		WrongAttempts:      entity.WrongAttempts,
		MaxWrongAttempts:   &entity.MaxWrongAttempts,
		RecoveryToken:      entity.RecoveryToken,
		MaxTokenAgeSeconds: &entity.MaxTokenAgeSeconds,
	})
	if err != nil {
		return nil, err
	}

	return domain, nil
}

func DomainToModel(domain models.Auth) entities.Auth {
	return entities.Auth{
		ID:                 domain.GetID(),
		IdentifierType:     IdentifierTypeFromDomain(domain.GetIdentifierType()),
		IdentifierValue:    domain.GetIdentifierValue(),
		Password:           domain.GetPassword(),
		EncryptToken:       domain.GetEncryptToken(),
		LastLoginAt:        domain.GetLastLoginAt(),
		WrongAttempts:      domain.GetWrongAttempts(),
		MaxWrongAttempts:   *domain.GetMaxWrongAttempts(),
		RecoveryToken:      domain.GetRecoveryToken(),
		MaxTokenAgeSeconds: *domain.GetMaxTokenAgeSeconds(),
	}
}