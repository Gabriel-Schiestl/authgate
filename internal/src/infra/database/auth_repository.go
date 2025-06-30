package database

import (
	"context"
	"fmt"

	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/models"
	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/repositories"
	"github.com/Gabriel-Schiestl/authgate/internal/src/infra/entities"
	"github.com/Gabriel-Schiestl/authgate/internal/src/infra/mappers"
	"github.com/Gabriel-Schiestl/go-clarch/v2/domain/exceptions"
	"gorm.io/gorm"
)

type authRepository struct{
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) repositories.IAuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) Save(ctx context.Context, auth models.Auth) (models.Auth, error) {
	authEntity := mappers.DomainToModel(auth)

	if err := r.db.WithContext(ctx).Save(&authEntity).Error; err != nil {
		return nil, err
	}

	return auth, nil
}

func (r *authRepository) GetByUserID(ctx context.Context, userID string) (models.Auth, error) {
    var authEntity entities.Auth

    if err := r.db.WithContext(ctx).Preload("UserInfo").
        Joins("JOIN user_info ON user_info.auth_id = auth.id").
        Where("user_info.user_id = ?", userID).First(&authEntity).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, exceptions.NewRepositoryNoDataFoundException(fmt.Sprintf("Auth not found for user ID: %s", userID))
        }
        return nil, err
    }

    auth, err := mappers.ModelToDomain(authEntity)
    if err != nil {
        return nil, err
    }

    return auth, nil
}

func (r *authRepository) GetByIdentifier(ctx context.Context, identifierType int, identifierValue string) (models.Auth, error) {
	var authEntity entities.Auth
	if err := r.db.WithContext(ctx).Preload("UserInfo").Where("identifier_type = ? AND identifier_value = ?", identifierType, identifierValue).First(&authEntity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exceptions.NewRepositoryNoDataFoundException(fmt.Sprintf("Auth not found for identifier type: %d and value: %s", identifierType, identifierValue))
		}
		return nil, err
	}

	auth, err := mappers.ModelToDomain(authEntity)
	if err != nil {
		return nil, err
	}

	return auth, nil
}

func (r *authRepository) Delete(ctx context.Context, userID string) error {
    var authEntity entities.Auth
    if err := r.db.WithContext(ctx).
        Joins("JOIN user_info ON user_info.auth_id = auth.id").
        Where("user_info.user_id = ?", userID).
        First(&authEntity).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return exceptions.NewRepositoryNoDataFoundException(fmt.Sprintf("Auth not found for user ID: %s", userID))
        }
        return err
    }

    if err := r.db.WithContext(ctx).Delete(&authEntity).Error; err != nil {
        return err
    }

    return nil
}