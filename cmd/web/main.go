package main

import (
	"github.com/VulpesFerrilata/auth/initialize"
	"github.com/VulpesFerrilata/auth/internal/pkg/micro/flags"
	common_flags "github.com/VulpesFerrilata/library/pkg/micro/flags"
	"github.com/asim/go-micro/plugins/server/http/v3"
	"github.com/asim/go-micro/v3"
	"github.com/kataras/iris/v12"
	"github.com/micro/cli/v2"
	"github.com/pkg/errors"
)

func main() {
	service := micro.NewService(
		micro.Server(
			http.NewServer(),
		),
		micro.Name("boardgame.auth.web"),
		micro.Version("latest"),
		micro.Flags(
			common_flags.NewSqlDialectFlag(),
			common_flags.NewSqlDsnFlag(),
			common_flags.NewTranslationFolderPathFlag(),
			flags.NewAccessTokenAlgFlag(),
			flags.NewAccessTokenSecretFlag(),
			flags.NewAccessTokenDurationFlag(),
			flags.NewRefreshTokenAlgFlag(),
			flags.NewRefreshTokenSecretFlag(),
			flags.NewRefreshTokenDurationFlag(),
		),
	)

	var cliCtx *cli.Context
	// Initialise service
	service.Init(
		micro.Action(func(ctx *cli.Context) error {
			cliCtx = ctx
			return nil
		}),
	)

	container := initialize.InitContainer(cliCtx)

	if err := container.Invoke(func(app *iris.Application) error {
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
