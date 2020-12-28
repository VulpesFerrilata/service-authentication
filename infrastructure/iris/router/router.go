package router

import (
	"database/sql"

	"github.com/VulpesFerrilata/auth/infrastructure/iris/controller"
	"github.com/VulpesFerrilata/library/pkg/middleware"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type Router interface {
	InitRoutes(app *iris.Application)
}

func NewRouter(authController controller.AuthController,
	transactionMiddleware *middleware.TransactionMiddleware,
	translatorMiddleware *middleware.TranslatorMiddleware,
	errorHandlerMiddleware *middleware.ErrorHandlerMiddleware) Router {
	return &router{
		authController:         authController,
		transactionMiddleware:  transactionMiddleware,
		translatorMiddleware:   translatorMiddleware,
		errorHandlerMiddleware: errorHandlerMiddleware,
	}
}

type router struct {
	authController         controller.AuthController
	transactionMiddleware  *middleware.TransactionMiddleware
	translatorMiddleware   *middleware.TranslatorMiddleware
	errorHandlerMiddleware *middleware.ErrorHandlerMiddleware
}

func (r router) InitRoutes(app *iris.Application) {
	apiRoot := app.Party("/api")
	apiRoot.Use(
		r.transactionMiddleware.ServeWithTxOptions(&sql.TxOptions{}),
		r.translatorMiddleware.Serve,
	)
	mvcApp := mvc.New(apiRoot.Party("/auth"))
	mvcApp.HandleError(r.errorHandlerMiddleware.ErrorHandler)
	mvcApp.Handle(r.authController)
}
