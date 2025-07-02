package usecases

import (
	"context"

	"github.com/Gabriel-Schiestl/authgate/internal/src/application/dtos"
	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/repositories"
	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/services"
	"github.com/Gabriel-Schiestl/go-clarch/v2/application/usecase"
	"github.com/Gabriel-Schiestl/go-clarch/v2/domain/exceptions"
)

type verifyTokenUsecase struct {
	authRepo repositories.IAuthRepository
	jwtService services.IJWTService
	encryptService services.IEncryptService
}

func NewVerifyTokenUsecase(authRepo repositories.IAuthRepository, jwtService services.IJWTService, encryptService services.IEncryptService) usecase.UseCaseWithProps[dtos.VerifyTokenDTO, *dtos.UserInfoDTO] {
	return &verifyTokenUsecase{
		authRepo: authRepo,
		jwtService: jwtService,
		encryptService: encryptService,
	}
}

func (luc verifyTokenUsecase) Execute(ctx context.Context, props dtos.VerifyTokenDTO) (*dtos.UserInfoDTO, error) {
	if props.AccessToken == "" {
		return nil, exceptions.NewBusinessException("access token is required")
	}

	token, _ := luc.encryptService.Decrypt(ctx, props.AccessToken)
	if token != "" {
		props.AccessToken = token
	}

	claims, err := luc.jwtService.ExtractClaims(ctx, props.AccessToken)
	if err != nil {
		return nil, exceptions.NewBusinessException("invalid access token")
	}

    auth, err := luc.authRepo.GetByUserID(ctx, claims["sub"].(string))
    if err != nil {
        return nil, err
    }

	return &dtos.UserInfoDTO{
		UserID: auth.GetUserInfo().GetUserID(),
		Name:   auth.GetUserInfo().GetName(),
		Roles:  auth.GetUserInfo().GetRoles(),
	}, nil
}