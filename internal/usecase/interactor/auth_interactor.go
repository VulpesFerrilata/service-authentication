package interactor

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/auth/internal/domain/service"
	"github.com/VulpesFerrilata/auth/internal/mapper"
	"github.com/VulpesFerrilata/auth/internal/pkg/app_error/authentication_error"
	"github.com/VulpesFerrilata/auth/internal/usecase/request"
	"github.com/VulpesFerrilata/auth/internal/usecase/response"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

	claim, err = a.claimService.Save(ctx, claim)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	standardClaim := mapper.NewClaimModelMapper(claim).ToStandardClaim()

	tokenResponse := new(response.TokenResponse)
	accessToken, err := a.tokenServiceFactory.GetTokenService(service.AccessToken).EncryptToken(ctx, standardClaim)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tokenResponse.AccessToken = accessToken

	refreshToken, err := a.tokenServiceFactory.GetTokenService(service.RefreshToken).EncryptToken(ctx, standardClaim)
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

	standardClaim, err := a.tokenServiceFactory.GetTokenService(service.AccessToken).DecryptToken(ctx, tokenRequest.Token)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userId, err := uuid.Parse(standardClaim.Subject)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	jti, err := uuid.Parse(standardClaim.Id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	claim, err := a.claimService.GetByUserId(ctx, userId)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if claim.GetJti() != jti {
		return nil, authentication_error.NewLoggedInByAnotherDeviceError()
	}

	return response.NewClaimResponse(claim), nil
}

func (a authInteractor) Refresh(ctx context.Context, tokenRequest *request.TokenRequest) (*response.TokenResponse, error) {
	if err := a.validate.StructCtx(ctx, tokenRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	standardClaim, err := a.tokenServiceFactory.GetTokenService(service.RefreshToken).DecryptToken(ctx, tokenRequest.Token)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userId, err := uuid.Parse(standardClaim.Subject)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	jti, err := uuid.Parse(standardClaim.Id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	claim, err := a.claimService.GetByUserId(ctx, userId)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if claim.GetJti() != jti {
		return nil, authentication_error.NewLoggedInByAnotherDeviceError()
	}

	tokenResponse := new(response.TokenResponse)
	accessToken, err := a.tokenServiceFactory.GetTokenService(service.AccessToken).EncryptToken(ctx, standardClaim)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tokenResponse.AccessToken = accessToken

	return tokenResponse, nil
}
