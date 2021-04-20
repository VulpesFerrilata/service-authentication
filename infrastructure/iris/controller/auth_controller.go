package controller

import (
	"github.com/VulpesFerrilata/auth/internal/usecase/interactor"
	"github.com/VulpesFerrilata/auth/internal/usecase/request"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/pkg/errors"
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

func (a authController) PostLogin(ctx iris.Context) interface{} {
	credentialRequest := new(request.CredentialRequest)

	if err := ctx.ReadJSON(credentialRequest); err != nil {
		return errors.WithStack(err)
	}

	tokenResponse, err := a.authInteractor.Login(ctx.Request().Context(), credentialRequest)
	if err != nil {
		return errors.WithStack(err)
	}

	return mvc.Response{
		Code:   iris.StatusOK,
		Object: tokenResponse,
	}
}

func (a authController) PostRefresh(ctx iris.Context) interface{} {
	tokenRequest := new(request.TokenRequest)

	if err := ctx.ReadJSON(tokenRequest); err != nil {
		return errors.WithStack(err)
	}

	tokenResponse, err := a.authInteractor.Refresh(ctx.Request().Context(), tokenRequest)
	if err != nil {
		return errors.WithStack(err)
	}

	return mvc.Response{
		Code:   iris.StatusOK,
		Object: tokenResponse,
	}
}
