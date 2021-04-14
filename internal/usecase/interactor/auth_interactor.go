package interactor

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/auth/internal/domain/service"
	"github.com/VulpesFerrilata/auth/internal/mapper"
	"github.com/VulpesFerrilata/auth/internal/usecase/request"
	"github.com/VulpesFerrilata/auth/internal/usecase/response"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type AuthInteractor interface {
	Login(ctx context.Context, credentialRequest *request.CredentialRequest) (*response.TokenResponse, error)
	Authenticate(ctx context.Context, tokenRequest *request.TokenRequest) (*response.ClaimResponse, error)
	Refresh(ctx context.Context, tokenRequest *request.TokenRequest) (*response.TokenResponse, error)
}

func NewAuthInteractor(validate *validator.Validate,
	userService service.UserService,
	claimService service.ClaimService,
	tokenServiceFactory service.TokenServiceFactory) AuthInteractor {
	return &authInteractor{
		validate:            validate,
		claimService:        claimService,
		tokenServiceFactory: tokenServiceFactory,
		userService:         userService,
	}
}

type authInteractor struct {
	validate            *validator.Validate
	claimService        service.ClaimService
	tokenServiceFactory service.TokenServiceFactory
	userService         service.UserService
}

func (a authInteractor) Login(ctx context.Context, credentialRequest *request.CredentialRequest) (*response.TokenResponse, error) {
	if err := a.validate.StructCtx(ctx, credentialRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	user, err := a.userService.GetByCredential(ctx, credentialRequest.Username, credentialRequest.Password)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	claim, err := model.NewClaim(user.GetId())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	claim, err = a.claimService.Create(ctx, claim)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tokenResponse := new(response.TokenResponse)
	accessStandardClaim := mapper.NewClaimModelMapper(claim).ToStandardClaim()
	accessToken, err := a.tokenServiceFactory.GetTokenService(service.AccessToken).EncryptToken(ctx, accessStandardClaim)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tokenResponse.AccessToken = accessToken

	refreshStandardClaim := mapper.NewClaimModelMapper(claim).ToStandardClaim()
	refreshToken, err := a.tokenServiceFactory.GetTokenService(service.RefreshToken).EncryptToken(ctx, refreshStandardClaim)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tokenResponse.RefreshToken = refreshToken

	return tokenResponse, nil
}

func (a authInteractor) Authenticate(ctx context.Context, tokenRequest *request.TokenRequest) (*response.ClaimResponse, error) {
	if err := a.validate.StructCtx(ctx, tokenRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	refreshStandardClaim, err := a.tokenServiceFactory.GetTokenService(service.RefreshToken).DecryptToken(ctx, tokenRequest.Token)
	if err != nil {
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Authenticate")
	}

	userId

	if err := a.claimService.GetByUserId(ctx, claim.GetUserId(), claim.GetJti()); err != nil {
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Authenticate")
	}

	return response.NewClaimResponse(claim), nil
}

func (a authInteractor) Refresh(ctx context.Context, tokenRequest *request.TokenRequest) (*response.TokenResponse, error) {
	if err := a.validate.StructCtx(ctx, tokenRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	claim, err := a.tokenService.DecryptToken(ctx, service.RefreshToken, tokenRequest.Token)
	if err != nil {
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Refresh")
	}

	if err := a.claimService.ValidateAuthenticate(ctx, claim.GetUserId(), claim.GetJti()); err != nil {
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Refresh")
	}

	tokenResponse := new(response.TokenResponse)
	accessToken, err := a.tokenService.EncryptToken(ctx, service.AccessToken, claim)
	if err != nil {
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Refresh")
	}
	tokenResponse.AccessToken = accessToken

	return tokenResponse, nil
}
