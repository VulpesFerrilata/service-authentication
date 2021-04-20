package main

import (
	"github.com/VulpesFerrilata/auth/initialize"
	grpcClient "github.com/VulpesFerrilata/go-micro/plugins/client/grpc"
	"github.com/asim/go-micro/plugins/server/http/v3"
	"github.com/asim/go-micro/v3"
	"github.com/kataras/iris/v12"
	"github.com/pkg/errors"
)

func main() {
	serviceName := micro.Name("boardgame.auth.web")
	serviceVersion := micro.Version("latest")
	serviceServer := micro.Server(
		http.NewServer(),
	)
	serviceClient := micro.Client(
		grpcClient.NewClient(),
	)

	container := initialize.InitContainer(serviceServer, serviceClient, serviceName, serviceVersion)

	if err := container.Invoke(func(service micro.Service, app *iris.Application) error {
		if err := app.Build(); err != nil {
			return errors.WithStack(err)
		}

		if err := micro.RegisterHandler(service.Server(), app); err != nil {
			return errors.WithStack(err)
		}

		// Run service
		if err := service.Run(); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		panic(err)
	}
}
