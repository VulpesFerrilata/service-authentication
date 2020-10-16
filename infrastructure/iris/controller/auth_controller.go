package controller

import (
	"github.com/VulpesFerrilata/auth/internal/usecase/interactor"
	"github.com/VulpesFerrilata/auth/internal/usecase/request"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type AuthController interface {
	PostLogin(ctx iris.Context) interface{}
}

func NewAuthController(authInteractor interactor.AuthInteractor) AuthController {
	return &authController{
		authInteractor: authInteractor,
	}
}

type authController struct {
	authInteractor interactor.AuthInteractor
}

func (ac authController) PostLogin(ctx iris.Context) interface{} {
	credentialRequest := new(request.CredentialRequest)

	if err := ctx.ReadJSON(credentialRequest); err != nil {
		return err
	}

	tokenResponse, err := ac.authInteractor.Login(ctx.Request().Context(), credentialRequest)
	if err != nil {
		return err
	}

	return mvc.Response{
		Code:   iris.StatusCreated,
		Object: tokenResponse,
	}
}
