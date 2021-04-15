package init

import (
	"github.com/VulpesFerrilata/auth/infrastructure/iris/controller"
	"github.com/VulpesFerrilata/auth/infrastructure/iris/router"
	"github.com/VulpesFerrilata/auth/infrastructure/iris/server"
	"github.com/VulpesFerrilata/auth/infrastructure/micro/handler"
	"github.com/VulpesFerrilata/auth/internal/domain/repository"
	"github.com/VulpesFerrilata/auth/internal/domain/service"
	"github.com/VulpesFerrilata/auth/internal/pkg/micro/flags"
	"github.com/VulpesFerrilata/auth/internal/usecase/interactor"
	"github.com/VulpesFerrilata/grpc/gateway"
	"github.com/VulpesFerrilata/library/init"
	common_flags "github.com/VulpesFerrilata/library/pkg/micro/flags"
	"github.com/VulpesFerrilata/library/pkg/middleware"
	ut "github.com/go-playground/universal-translator"
	"github.com/micro/cli/v2"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func InitContainer(ctx *cli.Context) *dig.Container {
	sqlDialect := common_flags.GetSqlDialect(ctx)
	sqlDsn := common_flags.GetSqlDsn(ctx)
	translationFolderPath := common_flags.GetTranslationFolderPath(ctx)
	accessTokenAlg := flags.GetAccessTokenAlg(ctx)
	accessTokenSecret := flags.GetAccessTokenSecret(ctx)
	accessTokenDuration := flags.GetAccessTokenDuration(ctx)
	refreshTokenAlg := flags.GetRefreshTokenAlg(ctx)
	refreshTokenSecret := flags.GetRefreshTokenSecret(ctx)
	refreshTokenDuration := flags.GetRefreshTokenDuration(ctx)

	container := dig.New()

	container.Provide(func() (*gorm.DB, error) {
		db, err := init.InitGorm(sqlDialect, sqlDsn)
		return db, errors.WithStack(err)
	})

	container.Provide(func() (*ut.UniversalTranslator, error) {
		utrans, err := init.InitTranslator(translationFolderPath)
		return utrans, errors.WithStack(err)
	})

	container.Provide(init.InitValidator)

	//--Middleware
	container.Provide(middleware.NewRecoverMiddleware)
	container.Provide(middleware.NewTransactionMiddleware)
	container.Provide(middleware.NewTranslatorMiddleware)
	container.Provide(middleware.NewErrorHandlerMiddleware)

	//--Repositories
	container.Provide(repository.NewClaimRepository)
	//--Services
	container.Provide(service.NewUserService)
	container.Provide(service.NewClaimService)
	container.Provide(func() service.TokenServiceFactory {
		accessTokenService := service.NewTokenService(accessTokenAlg, accessTokenSecret, accessTokenDuration)
		refreshTokenService := service.NewTokenService(refreshTokenAlg, refreshTokenSecret, refreshTokenDuration)
		return service.NewTokenServiceFactory(accessTokenService, refreshTokenService)
	})
	//--Usecases
	container.Provide(interactor.NewAuthInteractor)
	//--Gateways
	container.Provide(func() gateway.UserGateway {
		return gateway.NewUserGateway()
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
