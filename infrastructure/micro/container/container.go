package container

import (
	common_container "github.com/VulpesFerrilata/auth/infrastructure/container"
	"github.com/VulpesFerrilata/auth/infrastructure/micro/handler"
	"github.com/micro/cli/v2"
	"go.uber.org/dig"
)

func NewContainer(ctx *cli.Context) *dig.Container {
	container := common_container.NewContainer(ctx)

	//--Grpc
	container.Provide(handler.NewAuthHandler)

	return container
}
