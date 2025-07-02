package controller

import (
	"context"

	"github.com/Gabriel-Schiestl/authgate/internal/src/application/dtos"
	"github.com/Gabriel-Schiestl/go-clarch/v2/application/usecase"
)

type Controller struct {
	loginUsecase usecase.UseCaseWithProps[dtos.LoginDTO, *dtos.LoginResponseDTO]
	registerUsecase usecase.UseCaseWithProps[dtos.RegisterDTO, *dtos.RegisterResponseDTO]
	refreshUsecase usecase.UseCaseWithProps[dtos.RefreshTokenDTO, *dtos.RefreshTokenResponseDTO]
	verifyUsecase usecase.UseCaseWithProps[dtos.VerifyTokenDTO, *dtos.UserInfoDTO]
	deleteAuthUsecase usecase.UseCaseWithProps[string, *struct{}]
}

func NewController(
	loginUsecase usecase.UseCaseWithProps[dtos.LoginDTO, *dtos.LoginResponseDTO],
	registerUsecase usecase.UseCaseWithProps[dtos.RegisterDTO, *dtos.RegisterResponseDTO],
	refreshUsecase usecase.UseCaseWithProps[dtos.RefreshTokenDTO, *dtos.RefreshTokenResponseDTO],
	verifyUsecase usecase.UseCaseWithProps[dtos.VerifyTokenDTO, *dtos.UserInfoDTO],
	deleteAuthUsecase usecase.UseCaseWithProps[string, *struct{}],
) *Controller {
	controller := &Controller{
		loginUsecase: loginUsecase,
		registerUsecase: registerUsecase,
		refreshUsecase: refreshUsecase,
		verifyUsecase: verifyUsecase,
		deleteAuthUsecase: deleteAuthUsecase,
	}

	return controller
}

func (c *Controller) Login(ctx context.Context, dto dtos.LoginDTO) (*dtos.LoginResponseDTO, error) {
	response, err := usecase.ExecuteUseCaseWithProps(ctx, c.loginUsecase, dto)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Controller) Register(ctx context.Context, dto dtos.RegisterDTO) (*dtos.RegisterResponseDTO, error) {
	response, err := usecase.ExecuteUseCaseWithProps(ctx, c.registerUsecase, dto)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Controller) RefreshToken(ctx context.Context, dto dtos.RefreshTokenDTO) (*dtos.RefreshTokenResponseDTO, error) {
	response, err := usecase.ExecuteUseCaseWithProps(ctx, c.refreshUsecase, dto)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Controller) VerifyToken(ctx context.Context, dto dtos.VerifyTokenDTO) (*dtos.UserInfoDTO, error) {
	response, err := usecase.ExecuteUseCaseWithProps(ctx, c.verifyUsecase, dto)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Controller) DeleteAuth(ctx context.Context, userID string) error {
	_, err := usecase.ExecuteUseCaseWithProps(ctx, c.deleteAuthUsecase, userID)
	if err != nil {
		return err
	}

	return nil
}