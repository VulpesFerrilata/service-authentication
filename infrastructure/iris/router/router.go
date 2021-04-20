package router

import (
	"database/sql"

	"github.com/VulpesFerrilata/auth/infrastructure/iris/controller"
	"github.com/VulpesFerrilata/auth/internal/pkg/app_error/authentication_error"
	"github.com/VulpesFerrilata/library/pkg/middleware"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/pkg/errors"
)

type Router interface {
	InitRoutes(app *iris.Application)
}

func NewRouter(authController controller.AuthController,
	recoverMiddleware *middleware.RecoverMiddleware,
	transactionMiddleware *middleware.TransactionMiddleware,
	translatorMiddleware *middleware.TranslatorMiddleware,
	errorHandlerMiddleware *middleware.ErrorHandlerMiddleware) Router {
	return &router{
		authController:         authController,
		recoverMiddleware:      recoverMiddleware,
		transactionMiddleware:  transactionMiddleware,
		translatorMiddleware:   translatorMiddleware,
		errorHandlerMiddleware: errorHandlerMiddleware,
	}
}

type router struct {
	authController         controller.AuthController
	recoverMiddleware      *middleware.RecoverMiddleware
	transactionMiddleware  *middleware.TransactionMiddleware
	translatorMiddleware   *middleware.TranslatorMiddleware
	errorHandlerMiddleware *middleware.ErrorHandlerMiddleware
}

func (r router) InitRoutes(app *iris.Application) {
	apiRoot := app.Party("/api")
	apiRoot.Use(
		r.recoverMiddleware.Serve,
		r.transactionMiddleware.ServeWithTxOptions(&sql.TxOptions{}),
		r.translatorMiddleware.Serve,
	)
	mvcApp := mvc.New(apiRoot.Party("/auth"))
	mvcApp.HandleError(func(ctx iris.Context, err error) {
		if authenticationErr, ok := errors.Cause(err).(authentication_error.AuthenticationError); ok {
			err = authentication_error.NewAuthenticationErrors(authenticationErr)
		}

		r.errorHandlerMiddleware.ErrorHandler(ctx, err)
	})
	mvcApp.Handle(r.authController)
}
