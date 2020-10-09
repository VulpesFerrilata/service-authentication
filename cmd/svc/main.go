package main

import (
	"database/sql"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/server/grpc"

	"github.com/VulpesFerrilata/auth/infrastructure/container"
	"github.com/VulpesFerrilata/grpc/protoc/auth"
	"github.com/VulpesFerrilata/library/pkg/middleware"
	"github.com/micro/go-micro/v2"
)

func main() {
	container := container.NewContainer()

	if err := container.Invoke(func(authHandler auth.AuthHandler,
		transactionMiddleware *middleware.TransactionMiddleware,
		translatorMiddleware *middleware.TranslatorMiddleware,
		errorMiddleware *middleware.ErrorMiddleware) error {
		// New Service
		service := micro.NewService(
			micro.Server(
				grpc.NewServer(
					server.WrapHandler(errorMiddleware.HandlerWrapper),
					server.WrapHandler(transactionMiddleware.HandlerWrapperWithTxOptions(&sql.TxOptions{})),
					server.WrapHandler(translatorMiddleware.HandlerWrapper),
				),
			),
			micro.Name("boardgame.auth.svc"),
			micro.Version("latest"),
		)

		// Initialise service
		service.Init()

		// Register Handler
		if err := auth.RegisterAuthHandler(service.Server(), authHandler); err != nil {
			return err
		}

		// Run service
		return service.Run()
	}); err != nil {
		log.Fatal(err)
	}
}
