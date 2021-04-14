package container

import (
	common_container "github.com/VulpesFerrilata/auth/infrastructure/container"
	"github.com/VulpesFerrilata/auth/infrastructure/iris/controller"
	"github.com/VulpesFerrilata/auth/infrastructure/iris/router"
	"github.com/VulpesFerrilata/auth/infrastructure/iris/server"
	"github.com/micro/cli/v2"
	"go.uber.org/dig"
)

func NewContainer(ctx *cli.Context) *dig.Container {
	container := common_container.NewContainer(ctx)

	//--Controller
	container.Provide(controller.NewAuthController)
	//--Router
	container.Provide(router.NewRouter)
	//--Server
	container.Provide(server.NewServer)

	return container
}
