package interactor

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/datamodel"
	"github.com/VulpesFerrilata/auth/internal/domain/service"
	"github.com/VulpesFerrilata/auth/internal/usecase/request"
	"github.com/VulpesFerrilata/auth/internal/usecase/response"
	"github.com/VulpesFerrilata/grpc/protoc/user"
	"github.com/VulpesFerrilata/library/pkg/app_error"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

type AuthInteractor interface {
	Login(ctx context.Context, credentialRequest *request.CredentialRequest) (*response.TokenResponse, error)
	Authenticate(ctx context.Context, tokenRequest *request.TokenRequest) (*response.ClaimResponse, error)
	Refresh(ctx context.Context, tokenRequest *request.TokenRequest) (*response.TokenResponse, error)
}

func NewAuthInteractor(validate *validator.Validate,
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
	validate     *validator.Validate
	claimService service.ClaimService
	tokenService service.TokenService
	userService  user.UserService
}

func (ai authInteractor) Login(ctx context.Context, credentialRequest *request.CredentialRequest) (*response.TokenResponse, error) {
	if err := ai.validate.StructCtx(ctx, credentialRequest); err != nil {
		if fieldErrors, ok := errors.Cause(err).(validator.ValidationErrors); ok {
			err = app_error.NewValidationError(fieldErrors)
		}
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Login")
	}

	credentialRequestPb := new(user.CredentialRequest)
	credentialRequestPb.Username = credentialRequest.Username
	credentialRequestPb.Password = credentialRequest.Password
	userPb, err := ai.userService.GetUserByCredential(ctx, credentialRequestPb)
	if err != nil {
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Login")
	}

	claim, err := datamodel.NewClaim(int(userPb.GetID()))
	if err != nil {
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Login")
	}

	if err := ai.claimService.GetClaimRepository().InsertOrUpdate(ctx, claim); err != nil {
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Login")
	}

	tokenResponse := new(response.TokenResponse)
	accessToken, err := ai.tokenService.EncryptToken(ctx, service.AccessToken, claim)
	if err != nil {
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Login")
	}
	tokenResponse.AccessToken = accessToken

	refreshToken, err := ai.tokenService.EncryptToken(ctx, service.RefreshToken, claim)
	if err != nil {
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Login")
	}
	tokenResponse.RefreshToken = refreshToken

	return tokenResponse, nil
}

func (ai authInteractor) Authenticate(ctx context.Context, tokenRequest *request.TokenRequest) (*response.ClaimResponse, error) {
	if err := ai.validate.StructCtx(ctx, tokenRequest); err != nil {
		if fieldErrors, ok := errors.Cause(err).(validator.ValidationErrors); ok {
			err = app_error.NewValidationError(fieldErrors)
		}
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Authenticate")
	}

	claim, err := ai.tokenService.DecryptToken(ctx, service.AccessToken, tokenRequest.Token)
	if err != nil {
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Authenticate")
	}

	if err := ai.claimService.ValidateAuthenticate(ctx, claim.GetUserId(), claim.GetJti()); err != nil {
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Authenticate")
	}

	return response.NewClaimResponse(claim), nil
}

func (ai authInteractor) Refresh(ctx context.Context, tokenRequest *request.TokenRequest) (*response.TokenResponse, error) {
	if err := ai.validate.StructCtx(ctx, tokenRequest); err != nil {
		if fieldErrors, ok := errors.Cause(err).(validator.ValidationErrors); ok {
			err = app_error.NewValidationError(fieldErrors)
		}
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Refresh")
	}

	claim, err := ai.tokenService.DecryptToken(ctx, service.RefreshToken, tokenRequest.Token)
	if err != nil {
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Refresh")
	}

	if err := ai.claimService.ValidateAuthenticate(ctx, claim.GetUserId(), claim.GetJti()); err != nil {
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Refresh")
	}

	tokenResponse := new(response.TokenResponse)
	accessToken, err := ai.tokenService.EncryptToken(ctx, service.AccessToken, claim)
	if err != nil {
		return nil, errors.Wrap(err, "interactor.AuthInteractor.Refresh")
	}
	tokenResponse.AccessToken = accessToken

	return tokenResponse, nil
}
