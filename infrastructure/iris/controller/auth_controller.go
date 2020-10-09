package controller

import (
	"github.com/VulpesFerrilata/auth/infrastructure/iris/request"
	"github.com/VulpesFerrilata/auth/infrastructure/iris/response"
	"github.com/VulpesFerrilata/auth/internal/usecase/interactor"
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
	loginRequest := new(request.LoginRequest)

	if err := ctx.ReadJSON(loginRequest); err != nil {
		return err
	}

	tokenDTO, err := ac.authInteractor.Login(ctx.Request().Context(), loginRequest.ToCredentialRequestPb())
	if err != nil {
		return err
	}

	return mvc.Response{
		Code:   iris.StatusCreated,
		Object: response.NewTokenResponse(tokenDTO),
	}
}
