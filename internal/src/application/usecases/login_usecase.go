package usecases

import (
	"context"
	"fmt"

	"github.com/Gabriel-Schiestl/authgate/internal/src/application/dtos"
	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/repositories"
	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/services"
	"github.com/Gabriel-Schiestl/authgate/internal/src/utils"
	"github.com/Gabriel-Schiestl/go-clarch/v2/application/usecase"
	"github.com/Gabriel-Schiestl/go-clarch/v2/domain/exceptions"
)

type loginUsecase struct {
	authRepo repositories.IAuthRepository
	jwtService services.IJWTService
    encryptService services.IEncryptService
}

func NewLoginUsecase(authRepo repositories.IAuthRepository, jwtService services.IJWTService, encryptService services.IEncryptService) usecase.UseCaseWithProps[dtos.LoginDTO, *dtos.LoginResponseDTO] {
	return &loginUsecase{
		authRepo: authRepo,
		jwtService: jwtService,
        encryptService: encryptService,
	}
}

func (luc loginUsecase) Execute(ctx context.Context, props dtos.LoginDTO) (*dtos.LoginResponseDTO, error) {
    auth, err := luc.authRepo.GetByIdentifier(ctx, int(props.IdentifierType), props.IdentifierValue)
    if err != nil {
        return nil, err
    }

    if ok := utils.CheckPasswordHash(props.Password, auth.GetPassword()); !ok {
        return nil, exceptions.NewBusinessException("invalid credentials")
    }

    accessToken, err := luc.jwtService.GenerateToken(
        ctx,
        auth.GetUserInfo().GetUserID(), 
        auth.GetUserInfo().GetRoles(), 
        *auth.GetMaxTokenAgeSeconds(),
    )
    if err != nil {
        return nil, err
    }

    refreshTokenExpiry := 7 * 24 * 60 * 60
    refreshToken, err := luc.jwtService.GenerateRefreshToken(
        ctx,
        auth.GetUserInfo().GetUserID(),
        refreshTokenExpiry,
    )
    if err != nil {
        return nil, err
    }

    if auth.GetEncryptToken() {
        accessToken, err = luc.encryptService.Encrypt(ctx, *accessToken)
        if err != nil {
            fmt.Println("Failed to encrypt access token:", err)
            return nil, exceptions.NewBusinessException("failed to encrypt access token")
        }

        refreshToken, err = luc.encryptService.Encrypt(ctx, *refreshToken)
        if err != nil {
            return nil, exceptions.NewBusinessException("failed to encrypt refresh token")
        }
    }

    return &dtos.LoginResponseDTO{
        AccessToken:  *accessToken,
        RefreshToken: *refreshToken,
        UserInfo: dtos.UserInfoDTO{
            UserID: auth.GetUserInfo().GetUserID(),
            Name:   auth.GetUserInfo().GetName(),
            Roles:  auth.GetUserInfo().GetRoles(),
        },
    }, nil
}