package main

import (
	"context"
	"database/sql"

	"github.com/asim/go-micro/v3/server"
	"github.com/pkg/errors"

	"github.com/VulpesFerrilata/auth/initialize"
	"github.com/VulpesFerrilata/auth/internal/pkg/app_error/authentication_error"
	grpcClient "github.com/VulpesFerrilata/go-micro/plugins/client/grpc"
	grpcServer "github.com/VulpesFerrilata/go-micro/plugins/server/grpc"
	"github.com/VulpesFerrilata/grpc/protoc/auth"
	"github.com/VulpesFerrilata/library/pkg/middleware"
	"github.com/asim/go-micro/v3"
)

func main() {
	serviceName := micro.Name("boardgame.auth.svc")
	serviceVersion := micro.Version("latest")
	serviceServer := micro.Server(
		grpcServer.NewServer(),
	)
	serviceClient := micro.Client(
		grpcClient.NewClient(),
	)

	container := initialize.InitContainer(serviceServer, serviceClient, serviceName, serviceVersion)

	if err := container.Invoke(func(service micro.Service, authHandler auth.AuthHandler,
		recoverMiddleware *middleware.RecoverMiddleware,
		transactionMiddleware *middleware.TransactionMiddleware,
		translatorMiddleware *middleware.TranslatorMiddleware,
		errorHandlerMiddleware *middleware.ErrorHandlerMiddleware) error {

		service.Server().Init(
			server.WrapHandler(recoverMiddleware.HandlerWrapper),
			server.WrapHandler(transactionMiddleware.HandlerWrapperWithTxOptions(&sql.TxOptions{})),
			server.WrapHandler(translatorMiddleware.HandlerWrapper),
			server.WrapHandler(func(hf server.HandlerFunc) server.HandlerFunc {
				wrapper := func(ctx context.Context, request server.Request, response interface{}) error {
					err := hf(ctx, request, response)

					if authenticationErr, ok := errors.Cause(err).(authentication_error.AuthenticationError); ok {
						err = authentication_error.NewAuthenticationErrors(authenticationErr)
					}

					return err
				}

				return errorHandlerMiddleware.HandlerWrapper(wrapper)
			}),
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
