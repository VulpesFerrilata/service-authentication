package handler

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/usecase/input"
	"github.com/VulpesFerrilata/auth/internal/usecase/interactor"
	"github.com/VulpesFerrilata/grpc/protoc/auth"
	"github.com/pkg/errors"
)

func NewAuthHandler(authInteractor interactor.AuthInteractor) auth.AuthHandler {
	return &authHandler{
		authInteractor: authInteractor,
	}
}

type authHandler struct {
	authInteractor interactor.AuthInteractor
}

func (a authHandler) CreateUserCredential(ctx context.Context, userCredentialInputPb *auth.UserCredentialRequest, userCredentialResponsePb *auth.UserCredentialResponse) error {
	userCredentialInput := new(input.UserCredentialInput)
	userCredentialInput.UserID = userCredentialInputPb.GetUserID()
	userCredentialInput.Password = userCredentialInputPb.GetPassword()

	userCredentialOutput, err := a.authInteractor.CreateUserCredential(ctx, userCredentialInput)
	if err != nil {
		return errors.WithStack(err)
	}

	userCredentialResponsePb.ID = userCredentialOutput.ID

	return nil
}

func (a authHandler) Authenticate(ctx context.Context, tokenInputPb *auth.TokenRequest, claimResponsePb *auth.ClaimResponse) error {
	tokenInput := new(input.TokenInput)
	tokenInput.Token = tokenInputPb.GetToken()

	claimResponse, err := a.authInteractor.Authenticate(ctx, tokenInput)
	if err != nil {
		return errors.WithStack(err)
	}

	claimResponsePb.UserID = claimResponse.UserID

	return nil
}
