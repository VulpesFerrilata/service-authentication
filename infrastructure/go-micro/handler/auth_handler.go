package handler

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/usecase/interactor"
	"github.com/VulpesFerrilata/auth/internal/usecase/request"
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

func (ah authHandler) Authenticate(ctx context.Context, tokenRequestPb *auth.TokenRequest, claimResponsePb *auth.ClaimResponse) error {
	tokenRequest := new(request.TokenRequest)
	tokenRequest.Token = tokenRequestPb.Token

	claimResponse, err := ah.authInteractor.Authenticate(ctx, tokenRequest)
	if err != nil {
		return errors.Wrap(err, "handler.AuthHandler.Authenticate")
	}
	claimResponsePb.UserID = int64(claimResponse.UserID)

	return nil
}
