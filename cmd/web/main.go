package main

import (
	"log"

	"github.com/VulpesFerrilata/auth/infrastructure/container"
	"github.com/kataras/iris/v12"
	"github.com/micro/go-micro/v2/web"
)

func main() {
	container := container.NewContainer()

	if err := container.Invoke(func(app *iris.Application) error {
		// New Service
		service := web.NewService(
			web.Name("boardgame.auth.web"),
			web.Version("latest"),
		)

		// Initialise service
		if err := service.Init(); err != nil {
			return err
		}

		// Register Handler
		if err := app.Build(); err != nil {
			return err
		}
		service.Handle("/", app)

		// Run service
		return service.Run()
	}); err != nil {
		log.Fatal(err)
	}
}
