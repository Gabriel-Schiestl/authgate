package mappers

import (
	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/models"
	"github.com/Gabriel-Schiestl/authgate/internal/src/infra/entities"
)

func IdentifierTypeToDomain(it entities.IdentifierType) models.IdentifierType {
	return models.IdentifierType(it)
}

func IdentifierTypeFromDomain(domainType models.IdentifierType) entities.IdentifierType {
	return entities.IdentifierType(domainType)
}