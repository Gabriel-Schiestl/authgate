package module

import (
	"github.com/Gabriel-Schiestl/authgate/internal/src/application/usecases"
	"github.com/Gabriel-Schiestl/authgate/internal/src/controller"
	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/repositories"
	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/services"
	"github.com/Gabriel-Schiestl/authgate/internal/src/infra/adapters"
	"github.com/Gabriel-Schiestl/authgate/internal/src/infra/database"
	"github.com/Gabriel-Schiestl/authgate/internal/src/infra/database/connection"
	"github.com/Gabriel-Schiestl/authgate/internal/src/server"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func Module() fx.Option {
	return fx.Module(
		"authgate.module.app",
		fx.Provide(
			connection.SetupConfig,
		),
		fx.Provide(
			fx.Annotate(
				database.NewAuthRepository,
				fx.As(new(repositories.IAuthRepository)),
			),
			fx.Annotate(
				adapters.NewJWTService,
				fx.As(new(services.IJWTService)),
			),
			fx.Annotate(
				adapters.NewEncryptService,
				fx.As(new(services.IEncryptService)),
			),
			usecases.NewLoginUsecase,
			usecases.NewRegisterUsecase,
			usecases.NewVerifyTokenUsecase,
			usecases.NewRefreshTokenUsecase,
			usecases.NewDeleteAuthUsecase,
			controller.NewController,
		),
		fx.Invoke(server.NewAuthServiceServer),
		fx.Invoke(func(lc fx.Lifecycle, db *gorm.DB) {
			lc.Append(fx.StopHook(func() {
				sqlDb, err := db.DB()
				if err != nil {
					panic("failed to get database connection: " + err.Error())
				}
				
				if err := sqlDb.Close(); err != nil {
					panic("failed to close database connection: " + err.Error())
				}
			}))
		}),
	)
}