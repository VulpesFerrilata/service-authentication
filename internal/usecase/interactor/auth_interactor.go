package interactor

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/auth/internal/domain/service"
	"github.com/VulpesFerrilata/auth/internal/usecase/request"
	"github.com/VulpesFerrilata/auth/internal/usecase/response"
	"github.com/VulpesFerrilata/grpc/protoc/user"
	"github.com/VulpesFerrilata/library/pkg/validator"
)

type AuthInteractor interface {
	Login(ctx context.Context, credentialRequest *request.CredentialRequest) (*response.TokenResponse, error)
	Authenticate(ctx context.Context, tokenRequest *request.TokenRequest) (*response.ClaimResponse, error)
	Refresh(ctx context.Context, tokenRequest *request.TokenRequest) (*response.TokenResponse, error)
}

func NewAuthInteractor(validate validator.Validate,
	claimService service.ClaimService,
	tokenService service.TokenService,
	userService user.UserService) AuthInteractor {
	return &authInteractor{
		validate:     validate,
		claimService: claimService,
		tokenService: tokenService,
		userService:  userService,
	}
}

type authInteractor struct {
	validate     validator.Validate
	claimService service.ClaimService
	tokenService service.TokenService
	userService  user.UserService
}

func (ai authInteractor) Login(ctx context.Context, credentialRequest *request.CredentialRequest) (*response.TokenResponse, error) {
	if err := ai.validate.Struct(ctx, credentialRequest); err != nil {
		return nil, err
	}

	credentialRequestPb := credentialRequest.ToCredentialRequestPb()
	userPb, err := ai.userService.GetUserByCredential(ctx, credentialRequestPb)
	if err != nil {
		return nil, err
	}

	claim := new(model.Claim)
	claim.UserID = uint(userPb.ID)
	if err := ai.claimService.Create(ctx, claim); err != nil {
		return nil, err
	}

	tokenResponse := new(response.TokenResponse)
	accessToken, err := ai.tokenService.EncryptAccessToken(ctx, claim)
	if err != nil {
		return nil, err
	}
	tokenResponse.AccessToken = accessToken

	refreshToken, err := ai.tokenService.EncryptRefreshToken(ctx, claim)
	if err != nil {
		return nil, err
	}
	tokenResponse.RefreshToken = refreshToken

	return tokenResponse, nil
}

func (ai authInteractor) Authenticate(ctx context.Context, tokenRequest *request.TokenRequest) (*response.ClaimResponse, error) {
	if err := ai.validate.Struct(ctx, tokenRequest); err != nil {
		return nil, err
	}

	claim, err := ai.tokenService.DecryptAccessToken(ctx, tokenRequest.Token)
	if err != nil {
		return nil, err
	}

	if err := ai.claimService.ValidateAuthenticate(ctx, claim); err != nil {
		return nil, err
	}

	return response.NewClaimResponse(claim), nil
}

func (ai authInteractor) Refresh(ctx context.Context, tokenRequest *request.TokenRequest) (*response.TokenResponse, error) {
	if err := ai.validate.Struct(ctx, tokenRequest); err != nil {
		return nil, err
	}

	claim, err := ai.tokenService.DecryptRefreshToken(ctx, tokenRequest.Token)
	if err != nil {
		return nil, err
	}

	if err := ai.claimService.ValidateAuthenticate(ctx, claim); err != nil {
		return nil, err
	}

	tokenResponse := new(response.TokenResponse)
	accessToken, err := ai.tokenService.EncryptAccessToken(ctx, claim)
	if err != nil {
		return nil, err
	}
	tokenResponse.AccessToken = accessToken

	return tokenResponse, nil
}
