package handler

import (
	"context"

	"github.com/VulpesFerrilata/auth/infrastructure/go-micro/viewmodel"
	"github.com/VulpesFerrilata/auth/internal/usecase/interactor"
	"github.com/VulpesFerrilata/grpc/protoc/auth"
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
	tokenRequestVM := viewmodel.NewTokenRequest(tokenRequestPb)

	claimDTO, err := ah.authInteractor.Authenticate(ctx, tokenRequestVM.ToTokenForm())
	if err != nil {
		return err
	}

	claimResponseVM := viewmodel.NewClaimResponse(claimResponsePb)
	claimResponseVM.FromClaimDTO(claimDTO)
	return nil
}
