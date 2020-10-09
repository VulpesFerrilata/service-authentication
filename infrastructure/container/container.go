package container

import (
	"github.com/VulpesFerrilata/auth/infrastructure/go-micro/handler"
	"github.com/VulpesFerrilata/auth/infrastructure/iris/controller"
	"github.com/VulpesFerrilata/auth/infrastructure/iris/router"
	"github.com/VulpesFerrilata/auth/infrastructure/iris/server"
	"github.com/VulpesFerrilata/auth/internal/domain/repository"
	"github.com/VulpesFerrilata/auth/internal/domain/service"
	"github.com/VulpesFerrilata/auth/internal/usecase/adapter"
	"github.com/VulpesFerrilata/auth/internal/usecase/interactor"
	gateway "github.com/VulpesFerrilata/grpc/service"
	"github.com/VulpesFerrilata/library/config"
	"github.com/VulpesFerrilata/library/pkg/database"
	"github.com/VulpesFerrilata/library/pkg/db"
	"github.com/VulpesFerrilata/library/pkg/middleware"
	"github.com/VulpesFerrilata/library/pkg/translator"
	"github.com/VulpesFerrilata/library/pkg/validator"

	"go.uber.org/dig"
)

func NewContainer() *dig.Container {
	container := dig.New()

	//--Config
	container.Provide(config.NewConfig)
	container.Provide(config.NewJwtConfig)

	//--Domain
	container.Provide(repository.NewTokenRepository)
	container.Provide(service.NewAuthService)
	//--Usecase
	container.Provide(adapter.NewAuthAdapter)
	container.Provide(interactor.NewAuthInteractor)
	//--Gateways
	container.Provide(gateway.NewUserService)

	//--Utility
	container.Provide(database.NewGorm)
	container.Provide(db.NewDbContext)
	container.Provide(translator.NewTranslator)
	container.Provide(validator.NewValidate)

	//--Middleware
	container.Provide(middleware.NewTransactionMiddleware)
	container.Provide(middleware.NewTranslatorMiddleware)
	container.Provide(middleware.NewErrorMiddleware)

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
