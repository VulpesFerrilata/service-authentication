package controller

import (
	"github.com/VulpesFerrilata/auth/infrastructure/iris/request"
	"github.com/VulpesFerrilata/auth/infrastructure/iris/response"
	"github.com/VulpesFerrilata/auth/internal/usecase/input"
	"github.com/VulpesFerrilata/auth/internal/usecase/interactor"
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

	credentialInput := &input.CredentialInput{
		Username: credentialRequest.Username,
		Password: credentialRequest.Password,
	}

	tokenOutput, err := a.authInteractor.Login(ctx.Request().Context(), credentialInput)
	if err != nil {
		return errors.WithStack(err)
	}

	tokenResponse := &response.TokenResponse{
		AccessToken:  tokenOutput.AccessToken,
		RefreshToken: tokenOutput.RefreshToken,
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

	tokenInput := &input.TokenInput{
		Token: tokenRequest.Token,
	}

	tokenOutput, err := a.authInteractor.Refresh(ctx.Request().Context(), tokenInput)
	if err != nil {
		return errors.WithStack(err)
	}

	tokenResponse := &response.TokenResponse{
		AccessToken:  tokenOutput.AccessToken,
		RefreshToken: tokenOutput.RefreshToken,
	}
	return mvc.Response{
		Code:   iris.StatusOK,
		Object: tokenResponse,
	}
}
