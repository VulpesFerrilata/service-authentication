package interactor

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/auth/internal/domain/service"
	"github.com/VulpesFerrilata/auth/internal/pkg/app_error/authentication_error"
	"github.com/VulpesFerrilata/auth/internal/usecase/request"
	"github.com/VulpesFerrilata/auth/internal/usecase/response"
	"github.com/VulpesFerrilata/library/pkg/app_error"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthInteractor interface {
	CreateUserCredential(ctx context.Context, userCredentialRequest *request.UserCredentialRequest) error
	Login(ctx context.Context, credentialRequest *request.CredentialRequest) (*response.TokenResponse, error)
	Authenticate(ctx context.Context, tokenRequest *request.TokenRequest) (*response.ClaimResponse, error)
	Refresh(ctx context.Context, tokenRequest *request.TokenRequest) (*response.TokenResponse, error)
}

func NewAuthInteractor(validate *validator.Validate,
	userCredentialService service.UserCredentialService,
	claimService service.ClaimService,
	tokenServiceFactory service.TokenServiceFactory) AuthInteractor {
	return &authInteractor{
		validate:              validate,
		userCredentialService: userCredentialService,
		claimService:          claimService,
		tokenServiceFactory:   tokenServiceFactory,
	}
}

type authInteractor struct {
	validate              *validator.Validate
	userCredentialService service.UserCredentialService
	claimService          service.ClaimService
	tokenServiceFactory   service.TokenServiceFactory
}

func (a authInteractor) CreateUserCredential(ctx context.Context, userCredentialRequest *request.UserCredentialRequest) error {
	if err := a.validate.StructCtx(ctx, userCredentialRequest); err != nil {
		return errors.WithStack(err)
	}

	id, err := uuid.Parse(userCredentialRequest.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userCredentialRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.WithStack(err)
	}

	userCredential := model.NewUserCredential(
		id,
		userCredentialRequest.Username,
		hashPassword,
	)

	_, err = a.userCredentialService.Save(ctx, userCredential)
	return errors.WithStack(err)
}

func (a authInteractor) Login(ctx context.Context, credentialRequest *request.CredentialRequest) (*response.TokenResponse, error) {
	if err := a.validate.StructCtx(ctx, credentialRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	userCredential, err := a.userCredentialService.GetByUsername(ctx, credentialRequest.Username)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	jti, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	claim, err := a.claimService.GetById(ctx, userCredential.GetId())
	if err != nil && !app_error.IsRecordNotFoundError(errors.Cause(err)) {
		return nil, errors.WithStack(err)
	}
	if app_error.IsRecordNotFoundError(errors.Cause(err)) {
		claim = model.NewClaim(userCredential.GetId(), jti)
	} else {
		claim.SetJti(jti)
	}

	claim, err = a.claimService.Save(ctx, claim)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tokenResponse := new(response.TokenResponse)

	accessTokenStandardClaim := &jwt.StandardClaims{
		Id:      claim.GetJti().String(),
		Subject: claim.GetId().String(),
	}
	accessToken, err := a.tokenServiceFactory.GetTokenService(service.AccessToken).EncryptToken(ctx, accessTokenStandardClaim)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tokenResponse.AccessToken = accessToken

	refreshTokenStandardClaim := &jwt.StandardClaims{
		Id:      claim.GetJti().String(),
		Subject: claim.GetId().String(),
	}
	refreshToken, err := a.tokenServiceFactory.GetTokenService(service.RefreshToken).EncryptToken(ctx, refreshTokenStandardClaim)
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

	accessTokenStandardClaim, err := a.tokenServiceFactory.GetTokenService(service.AccessToken).DecryptToken(ctx, tokenRequest.Token)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	id, err := uuid.Parse(accessTokenStandardClaim.Subject)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	jti, err := uuid.Parse(accessTokenStandardClaim.Id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	claim, err := a.claimService.GetById(ctx, id)
	if err != nil {
		if app_error.IsRecordNotFoundError(errors.Cause(err)) {
			return nil, authentication_error.NewAuthenticationRevokedError()
		}
		return nil, errors.WithStack(err)
	}

	if claim.GetJti().String() == jti.String() {
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

	claim, err := a.claimService.GetById(ctx, userId)
	if err != nil {
		if app_error.IsRecordNotFoundError(errors.Cause(err)) {
			return nil, authentication_error.NewAuthenticationRevokedError()
		}
		return nil, errors.WithStack(err)
	}

	if claim.GetJti().String() == jti.String() {
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
