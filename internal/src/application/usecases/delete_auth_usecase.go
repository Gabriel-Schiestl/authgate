package usecases

import (
	"context"

	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/repositories"
	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/services"
	"github.com/Gabriel-Schiestl/go-clarch/v2/application/usecase"
	"github.com/Gabriel-Schiestl/go-clarch/v2/domain/exceptions"
)

type deleteAuthUsecase struct {
	authRepo repositories.IAuthRepository
	jwtService services.IJWTService
}

func NewDeleteAuthUsecase(authRepo repositories.IAuthRepository, jwtService services.IJWTService) usecase.UseCaseWithProps[string, *struct{}] {
	return &deleteAuthUsecase{
		authRepo: authRepo,
		jwtService: jwtService,
	}
}

func (luc deleteAuthUsecase) Execute(ctx context.Context, userID string) (*struct{}, error) {
	if userID == "" {
		return nil, exceptions.NewBusinessException("access token is required")
	}

    err := luc.authRepo.Delete(ctx, userID)
    if err != nil {
        return nil, err
    }

	return &struct{}{}, nil
}