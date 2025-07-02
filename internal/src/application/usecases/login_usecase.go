package usecases

import (
	"context"

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
}

func NewLoginUsecase(authRepo repositories.IAuthRepository, jwtService services.IJWTService) usecase.UseCaseWithProps[dtos.LoginDTO, *dtos.LoginResponseDTO] {
	return &loginUsecase{
		authRepo: authRepo,
		jwtService: jwtService,
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