package initialize

import (
	"github.com/VulpesFerrilata/auth/infrastructure/iris/controller"
	"github.com/VulpesFerrilata/auth/infrastructure/iris/router"
	"github.com/VulpesFerrilata/auth/infrastructure/iris/server"
	"github.com/VulpesFerrilata/auth/infrastructure/micro/handler"
	"github.com/VulpesFerrilata/auth/internal/domain/repository"
	"github.com/VulpesFerrilata/auth/internal/domain/service"
	"github.com/VulpesFerrilata/auth/internal/pkg/micro/flags"
	"github.com/VulpesFerrilata/auth/internal/usecase/interactor"
	"github.com/VulpesFerrilata/grpc/protoc/user"
	"github.com/VulpesFerrilata/library/initialize"
	common_flags "github.com/VulpesFerrilata/library/pkg/micro/flags"
	"github.com/VulpesFerrilata/library/pkg/middleware"
	"github.com/asim/go-micro/v3"
	ut "github.com/go-playground/universal-translator"
	"github.com/micro/cli/v2"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func InitContainer(opts ...micro.Option) *dig.Container {
	container := dig.New()

	container.Provide(func() micro.Service {
		serviceFlags := micro.Flags(
			common_flags.NewSqlDialectFlag(),
			common_flags.NewSqlDsnFlag(),
			common_flags.NewTranslationFolderPathFlag(),
			flags.NewAccessTokenAlgFlag(),
			flags.NewAccessTokenSecretFlag(),
			flags.NewAccessTokenDurationFlag(),
			flags.NewRefreshTokenAlgFlag(),
			flags.NewRefreshTokenSecretFlag(),
			flags.NewRefreshTokenDurationFlag(),
		)

		opts = append(opts, serviceFlags)

		return micro.NewService(opts...)
	})

	container.Provide(func(service micro.Service) *cli.Context {
		var ctx *cli.Context

		service.Init(
			micro.Action(func(c *cli.Context) error {
				ctx = c
				return nil
			}),
		)

		return ctx
	})

	container.Provide(func(ctx *cli.Context) (*gorm.DB, error) {
		sqlDialect := common_flags.GetSqlDialect(ctx)
		sqlDsn := common_flags.GetSqlDsn(ctx)

		db, err := initialize.InitGorm(sqlDialect, sqlDsn)
		return db, errors.WithStack(err)
	})

	container.Provide(func(ctx *cli.Context) (*ut.UniversalTranslator, error) {
		translationFolderPath := common_flags.GetTranslationFolderPath(ctx)

		utrans, err := initialize.InitTranslator(translationFolderPath)
		return utrans, errors.WithStack(err)
	})

	container.Provide(initialize.InitValidator)

	//--Middleware
	container.Provide(middleware.NewRecoverMiddleware)
	container.Provide(middleware.NewTransactionMiddleware)
	container.Provide(middleware.NewTranslatorMiddleware)
	container.Provide(middleware.NewErrorHandlerMiddleware)

	//--Repositories
	container.Provide(repository.NewClaimRepository)
	//--Services
	container.Provide(service.NewClaimService)
	container.Provide(func(ctx *cli.Context) service.TokenServiceFactory {
		accessTokenAlg := flags.GetAccessTokenAlg(ctx)
		accessTokenSecret := flags.GetAccessTokenSecret(ctx)
		accessTokenDuration := flags.GetAccessTokenDuration(ctx)
		refreshTokenAlg := flags.GetRefreshTokenAlg(ctx)
		refreshTokenSecret := flags.GetRefreshTokenSecret(ctx)
		refreshTokenDuration := flags.GetRefreshTokenDuration(ctx)

		accessTokenService := service.NewTokenService(accessTokenAlg, accessTokenSecret, accessTokenDuration)
		refreshTokenService := service.NewTokenService(refreshTokenAlg, refreshTokenSecret, refreshTokenDuration)

		return service.NewTokenServiceFactory(accessTokenService, refreshTokenService)
	})
	//--Usecases
	container.Provide(interactor.NewAuthInteractor)
	//--Gateways
	container.Provide(func(service micro.Service) user.UserService {
		return user.NewUserService("boardgame.user.svc", service.Client())
	})

	//--Controller
	container.Provide(controller.NewAuthController)
	//--Router
	container.Provide(router.NewRouter)
	//--Server
	container.Provide(server.NewServer)

	//--Grpc
	container.Provide(handler.NewAuthHandler)

	return container
}
