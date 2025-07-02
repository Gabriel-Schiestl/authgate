package usecases

import (
	"context"

	"github.com/Gabriel-Schiestl/authgate/internal/src/application/dtos"
	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/repositories"
	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/services"
	"github.com/Gabriel-Schiestl/go-clarch/v2/application/usecase"
	"github.com/Gabriel-Schiestl/go-clarch/v2/domain/exceptions"
)

type refreshTokenUsecase struct {
	authRepo repositories.IAuthRepository
	jwtService services.IJWTService
	encryptService services.IEncryptService
}

func NewRefreshTokenUsecase(authRepo repositories.IAuthRepository, jwtService services.IJWTService, encryptService services.IEncryptService) usecase.UseCaseWithProps[dtos.RefreshTokenDTO, *dtos.RefreshTokenResponseDTO] {
	return &refreshTokenUsecase{
		authRepo: authRepo,
		jwtService: jwtService,
		encryptService: encryptService,
	}
}

func (luc refreshTokenUsecase) Execute(ctx context.Context, props dtos.RefreshTokenDTO) (*dtos.RefreshTokenResponseDTO, error) {
	if props.RefreshToken == "" {
		return nil, exceptions.NewBusinessException("refresh token is required")
	}

	token, _ := luc.encryptService.Decrypt(ctx, props.RefreshToken)
	if token != "" {
		props.RefreshToken = token
	}

	claims, err := luc.jwtService.ExtractRefreshClaims(ctx, props.RefreshToken)
	if err != nil {
		return nil, exceptions.NewBusinessException("invalid refresh token")
	}

    auth, err := luc.authRepo.GetByUserID(ctx, claims["sub"].(string))
    if err != nil {
        return nil, err
    }

	newAccessToken, err := luc.jwtService.GenerateToken(
		ctx,
		auth.GetUserInfo().GetUserID(),
		auth.GetUserInfo().GetRoles(),
		*auth.GetMaxTokenAgeSeconds(),
	)
	if err != nil {
		return nil, err
	}

	if auth.GetEncryptToken() {
		newAccessToken, err = luc.encryptService.Encrypt(ctx, *newAccessToken)
		if err != nil {
			return nil, exceptions.NewBusinessException("failed to encrypt access token")
		}
	}

	return &dtos.RefreshTokenResponseDTO{
		AccessToken: *newAccessToken,
		UserInfo: dtos.UserInfoDTO{
			UserID: auth.GetUserInfo().GetUserID(),
			Name:   auth.GetUserInfo().GetName(),
			Roles:  auth.GetUserInfo().GetRoles(),
		},
	}, nil
}