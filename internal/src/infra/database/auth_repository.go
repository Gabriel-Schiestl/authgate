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

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) repositories.IAuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) Save(ctx context.Context, auth models.Auth) (models.Auth, error) {
	authEntity := mappers.DomainToModel(auth)

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("UserInfo").Create(&authEntity).Error; err != nil {
			return fmt.Errorf("failed to save auth: %w", err)
		}

		if authEntity.UserInfo.UserID != "" {
			authEntity.UserInfo.AuthID = authEntity.ID
			if err := tx.Create(&authEntity.UserInfo).Error; err != nil {
				return fmt.Errorf("failed to save user info: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).
        Preload("UserInfo").
        Where("id = ?", authEntity.ID).
        First(&authEntity).Error; err != nil {
        return nil, fmt.Errorf("failed to reload saved auth: %w", err)
    }

	savedAuth, err := mappers.ModelToDomain(authEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to convert saved entity to domain: %w", err)
	}

	return savedAuth, nil
}

func (r *authRepository) GetByUserID(ctx context.Context, userID string) (models.Auth, error) {
	var authEntity entities.Auth

	if err := r.db.WithContext(ctx).
		Preload("UserInfo").
		Joins("JOIN user_infos ON user_infos.auth_id = auths.id").
		Where("user_infos.user_id = ?", userID).
		First(&authEntity).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, exceptions.NewRepositoryNoDataFoundException(
				fmt.Sprintf("Auth not found for user ID: %s", userID))
		}
		return nil, fmt.Errorf("database error in GetByUserID: %w", err)
	}

	auth, err := mappers.ModelToDomain(authEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to convert entity to domain: %w", err)
	}

	return auth, nil
}

func (r *authRepository) GetByIdentifier(ctx context.Context, identifierType int, identifierValue string) (models.Auth, error) {
	var authEntity entities.Auth

	if err := r.db.WithContext(ctx).
		Preload("UserInfo").
		Where("identifier_type = ? AND identifier_value = ?", identifierType, identifierValue).
		First(&authEntity).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, exceptions.NewRepositoryNoDataFoundException(
				fmt.Sprintf("Auth not found for identifier type: %d and value: %s", identifierType, identifierValue))
		}
		return nil, fmt.Errorf("database error in GetByIdentifier: %w", err)
	}

	auth, err := mappers.ModelToDomain(authEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to convert entity to domain: %w", err)
	}

	return auth, nil
}

func (r *authRepository) Delete(ctx context.Context, userID string) error {
	var authEntity entities.Auth

	if err := r.db.WithContext(ctx).
		Preload("UserInfo").
		Joins("JOIN user_infos ON user_infos.auth_id = auths.id").
		Where("user_infos.user_id = ?", userID).
		First(&authEntity).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return exceptions.NewRepositoryNoDataFoundException(
				fmt.Sprintf("Auth not found for user ID: %s", userID))
		}
		return fmt.Errorf("database error finding auth to delete: %w", err)
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if authEntity.UserInfo.AuthID != "" {
			if err := tx.Delete(&authEntity.UserInfo).Error; err != nil {
				return fmt.Errorf("failed to delete user info: %w", err)
			}
		}

		// Deletar Auth
		if err := tx.Delete(&authEntity).Error; err != nil {
			return fmt.Errorf("failed to delete auth: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed during delete: %w", err)
	}

	return nil
}