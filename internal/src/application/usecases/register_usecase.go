package usecases

import (
	"context"
	"sync"

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

func NewRegisterUsecase(authRepo repositories.IAuthRepository) usecase.UseCaseWithProps[dtos.RegisterDTO, *dtos.RegisterResponseDTO] {
	return &registerUsecase{
		authRepo: authRepo,
	}
}

func (luc registerUsecase) Execute(ctx context.Context, props dtos.RegisterDTO) (*dtos.RegisterResponseDTO, error) {
    if  props.IdentifierValue == "" || props.Password == "" {
		return nil, exceptions.NewBusinessException("identifier type, identifier value, and password are required")
	}

	userExistsChan := make(chan models.Auth, 2)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		userExists, _ := luc.authRepo.GetByUserID(ctx, props.UserInfo.UserID)
		if userExists != nil {
			userExistsChan <- userExists
		}
	}()

	go func() {
		defer wg.Done()
		userExists, _ := luc.authRepo.GetByIdentifier(ctx, int(props.IdentifierType), props.IdentifierValue)
		if userExists != nil {
			userExistsChan <- userExists
		}
	}()

	go func() {
		wg.Wait()
		close(userExistsChan)
	}()
	
	for userExists := range userExistsChan {
		if userExists != nil {
			return nil, exceptions.NewBusinessException("auth already exists")
		}
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
		IdentifierType: auth.GetIdentifierType().String(),
		IdentifierValue: auth.GetIdentifierValue(),
		UserInfo: dtos.UserInfoDTO{
			UserID: auth.GetUserInfo().GetUserID(),
			Name:   auth.GetUserInfo().GetName(),
			Roles:  auth.GetUserInfo().GetRoles(),
		},
	}, nil
}