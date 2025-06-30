package repositories

import (
	"context"

	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/models"
)

type IAuthRepository interface {
	Save(ctx context.Context, auth models.Auth) (models.Auth, error)
	GetByUserID(ctx context.Context, userID string) (models.Auth, error)
	GetByIdentifier(ctx context.Context, identifierType int, identifierValue string) (models.Auth, error)
	Delete(ctx context.Context, userID string) error
	Update(ctx context.Context, auth models.Auth) (models.Auth, error)
}