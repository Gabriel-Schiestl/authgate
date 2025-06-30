package repositories

import "github.com/Gabriel-Schiestl/authgate/internal/src/domain/models"

type IAuthRepository interface {
	Save(auth models.Auth) (models.Auth, error)
	GetByUserID(userID string) (models.Auth, error)
	GetByIdentifier(identifierValue string) (models.Auth, error)
	Delete(userID string) error
	Update(auth models.Auth) (models.Auth, error)
}