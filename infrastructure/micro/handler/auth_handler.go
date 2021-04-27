package handler

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/usecase/interactor"
	"github.com/VulpesFerrilata/auth/internal/usecase/request"
	"github.com/VulpesFerrilata/grpc/protoc/auth"
	"github.com/golang/protobuf/ptypes/empty"
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

func (a authHandler) CreateUserCredential(ctx context.Context, userCredentialRequestPb *auth.UserCredentialRequest, emptyResponsePb *empty.Empty) error {
	userCredentialRequest := new(request.UserCredentialRequest)
	userCredentialRequest.ID = userCredentialRequestPb.GetID()
	userCredentialRequest.Username = userCredentialRequestPb.GetUsername()
	userCredentialRequest.Password = userCredentialRequestPb.GetPassword()

	err := a.authInteractor.CreateUserCredential(ctx, userCredentialRequest)
	return errors.WithStack(err)
}

func (a authHandler) Authenticate(ctx context.Context, tokenRequestPb *auth.TokenRequest, claimResponsePb *auth.ClaimResponse) error {
	tokenRequest := new(request.TokenRequest)
	tokenRequest.Token = tokenRequestPb.GetToken()

	claimResponse, err := a.authInteractor.Authenticate(ctx, tokenRequest)
	if err != nil {
		return errors.WithStack(err)
	}

	claimResponsePb.ID = claimResponse.ID

	return nil
}
