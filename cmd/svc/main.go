package main

import (
	"database/sql"

	"github.com/asim/go-micro/v3/server"
	"github.com/micro/cli/v2"
	"github.com/pkg/errors"

	"github.com/VulpesFerrilata/auth/infrastructure/micro/container"
	"github.com/VulpesFerrilata/auth/internal/pkg/micro/flags"
	"github.com/VulpesFerrilata/grpc/protoc/auth"
	"github.com/VulpesFerrilata/library/pkg/middleware"
	"github.com/asim/go-micro/v3"
)

func main() {
	service := micro.NewService(
		micro.Name("boardgame.auth.svc"),
		micro.Version("latest"),
		micro.Flags(
			flags.DefaultFlags...,
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

	container := container.NewContainer(cliCtx)

	if err := container.Invoke(func(authHandler auth.AuthHandler,
		recoverMiddleware *middleware.RecoverMiddleware,
		transactionMiddleware *middleware.TransactionMiddleware,
		translatorMiddleware *middleware.TranslatorMiddleware,
		errorHandlerMiddleware *middleware.ErrorHandlerMiddleware) error {
		// New Service
		service := micro.NewService(
			micro.Server(
				server.NewServer(
					server.WrapHandler(recoverMiddleware.HandlerWrapper),
					server.WrapHandler(errorHandlerMiddleware.HandlerWrapper),
					server.WrapHandler(translatorMiddleware.HandlerWrapper),
					server.WrapHandler(transactionMiddleware.HandlerWrapperWithTxOptions(&sql.TxOptions{})),
				),
			),
		)

		// Register Handler
		if err := auth.RegisterAuthHandler(service.Server(), authHandler); err != nil {
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
