package usecases

import (
	"context"

	"github.com/Gabriel-Schiestl/authgate/internal/src/application/dtos"
	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/models"
	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/repositories"
	"github.com/Gabriel-Schiestl/authgate/internal/src/utils"
	"github.com/Gabriel-Schiestl/go-clarch/v2/application/usecase"
	"github.com/Gabriel-Schiestl/go-clarch/v2/domain/exceptions"
)

type registerUsecase struct {
	authRepo repositories.IAuthRepository
}

type checkResult struct {
	exists bool
	err    error
}

func NewRegisterUsecase(authRepo repositories.IAuthRepository) usecase.UseCaseWithProps[dtos.RegisterDTO, *dtos.RegisterResponseDTO] {
	return &registerUsecase{
		authRepo: authRepo,
	}
}

func (luc registerUsecase) Execute(ctx context.Context, props dtos.RegisterDTO) (*dtos.RegisterResponseDTO, error) {
    if  props.IdentifierValue == "" || props.Password == "" {
		return nil, exceptions.NewBusinessException("identifier type, identifier value, and password are required")
	}

	userIDChan := make(chan checkResult, 1)
    identifierChan := make(chan checkResult, 1)

    // Verificar UserID em paralelo
    go func() {
        auth, err := luc.authRepo.GetByUserID(ctx, props.UserInfo.UserID)
        userIDChan <- checkResult{
            exists: auth != nil && err == nil,
            err:    err,
        }
    }()

    go func() {
        auth, err := luc.authRepo.GetByIdentifier(ctx, int(props.IdentifierType), props.IdentifierValue)
        identifierChan <- checkResult{
            exists: auth != nil && err == nil,
            err:    err,
        }
    }()

    userIDResult := <-userIDChan
    identifierResult := <-identifierChan

    if userIDResult.exists {
        return nil, exceptions.NewBusinessException("user with this ID already exists")
    }
    if identifierResult.exists {
        return nil, exceptions.NewBusinessException("user with this identifier already exists")
    }

	userInfo, err := models.NewUserInfo(models.UserInfoProps{
		UserID: props.UserInfo.UserID,
		Name:   props.UserInfo.Name,
		Roles:  props.UserInfo.Roles,
	})
	if err != nil {
		return nil, err
	}

	hashedPassword, er := utils.HashPassword(props.Password)
	if er != nil {
		return nil, exceptions.NewBusinessException("failed to hash password")
	}

	auth, err := models.NewAuth(models.AuthProps{
		IdentifierType:     models.IdentifierType(props.IdentifierType),
		IdentifierValue:    props.IdentifierValue,
		Password:           hashedPassword,
		UserInfo:           userInfo,
		EncryptToken:       props.EncryptToken,
		MaxTokenAgeSeconds: props.MaxTokenAgeSeconds,
		MaxWrongAttempts:   props.MaxWrongAttempts,
	})
	if err != nil {
		return nil, err
	}

	if _, err := luc.authRepo.Save(ctx, auth); err != nil {
		return nil, err
	}

	return &dtos.RegisterResponseDTO{
		IdentifierType: auth.GetIdentifierType(),
		IdentifierValue: auth.GetIdentifierValue(),
		UserInfo: dtos.UserInfoDTO{
			UserID: auth.GetUserInfo().GetUserID(),
			Name:   auth.GetUserInfo().GetName(),
			Roles:  auth.GetUserInfo().GetRoles(),
		},
	}, nil
}