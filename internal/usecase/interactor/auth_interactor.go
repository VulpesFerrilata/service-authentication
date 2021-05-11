package interactor

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/service"
	"github.com/VulpesFerrilata/auth/internal/pkg/app_error/detail_error"
	"github.com/VulpesFerrilata/auth/internal/usecase/input"
	"github.com/VulpesFerrilata/auth/internal/usecase/output"
	"github.com/VulpesFerrilata/library/pkg/app_error"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthInteractor interface {
	CreateUserCredential(ctx context.Context, userCredentialInput *input.UserCredentialInput) (*output.UserCredentialOutput, error)
	Login(ctx context.Context, credentialInput *input.CredentialInput) (*output.TokenOutput, error)
	Authenticate(ctx context.Context, tokenInput *input.TokenInput) (*output.ClaimOutput, error)
	Refresh(ctx context.Context, tokenInput *input.TokenInput) (*output.TokenOutput, error)
	Revoke(ctx context.Context, tokenInput *input.TokenInput) error
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

func (a authInteractor) CreateUserCredential(ctx context.Context, userCredentialInput *input.UserCredentialInput) (*output.UserCredentialOutput, error) {
	if err := a.validate.StructCtx(ctx, userCredentialInput); err != nil {
		return nil, errors.WithStack(err)
	}

	userId, err := uuid.Parse(userCredentialInput.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userCredentialInput.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userCredential, err := a.userCredentialService.NewUserCredential(ctx, userId, userCredentialInput.Username, hashPassword)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userCredential, err = a.userCredentialService.Save(ctx, userCredential)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userCredentialOutput := &output.UserCredentialOutput{
		ID:       userCredential.GetId().String(),
		Username: userCredential.GetUsername(),
	}
	return userCredentialOutput, nil
}

func (a authInteractor) Login(ctx context.Context, credentialInput *input.CredentialInput) (*output.TokenOutput, error) {
	authenticationErrs := app_error.NewAuthenticationErrors()

	if err := a.validate.StructCtx(ctx, credentialInput); err != nil {
		return nil, errors.WithStack(err)
	}

	userCredential, err := a.userCredentialService.GetByUsername(ctx, credentialInput.Username)
	if err != nil {
		detailErr := detail_error.NewInvalidPasswordError()
		authenticationErrs.AddDetailError(detailErr)
		return nil, authenticationErrs
	}

	if err := bcrypt.CompareHashAndPassword(userCredential.GetHashPassword(), []byte(credentialInput.Password)); err != nil {
		return nil, errors.WithStack(err)
	}

	jti, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	claim, err := a.claimService.GetById(ctx, userCredential.GetId())
	if err != nil && !app_error.IsRecordNotFoundError(err) {
		return nil, errors.WithStack(err)
	}
	if app_error.IsRecordNotFoundError(err) {
		claim, err = a.claimService.NewClaim(ctx, userCredential.GetId(), jti)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	claim.SetJti(jti)

	claim, err = a.claimService.Save(ctx, claim)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tokenOutput := new(output.TokenOutput)

	accessTokenStandardClaim := &jwt.StandardClaims{
		Id:      claim.GetJti().String(),
		Subject: claim.GetUserID().String(),
	}
	accessToken, err := a.tokenServiceFactory.GetTokenService(service.AccessToken).EncryptToken(ctx, accessTokenStandardClaim)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tokenOutput.AccessToken = accessToken

	refreshTokenStandardClaim := &jwt.StandardClaims{
		Id:      claim.GetJti().String(),
		Subject: claim.GetUserID().String(),
	}
	refreshToken, err := a.tokenServiceFactory.GetTokenService(service.RefreshToken).EncryptToken(ctx, refreshTokenStandardClaim)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tokenOutput.RefreshToken = refreshToken

	return tokenOutput, nil
}

func (a authInteractor) Authenticate(ctx context.Context, tokenInput *input.TokenInput) (*output.ClaimOutput, error) {
	authenticationErrors := app_error.NewAuthenticationErrors()

	if err := a.validate.StructCtx(ctx, tokenInput); err != nil {
		return nil, errors.WithStack(err)
	}

	accessTokenStandardClaim, err := a.tokenServiceFactory.GetTokenService(service.AccessToken).DecryptToken(ctx, tokenInput.Token)
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
			detailErr := detail_error.NewTokenRevokedError()
			authenticationErrors.AddDetailError(detailErr)
			return nil, authenticationErrors
		}
		return nil, errors.WithStack(err)
	}

	if claim.GetJti().String() == jti.String() {
		detailErr := detail_error.NewTokenRevokedError()
		authenticationErrors.AddDetailError(detailErr)
		return nil, authenticationErrors
	}

	claimOutput := &output.ClaimOutput{
		UserID: claim.GetUserID().String(),
	}
	return claimOutput, nil
}

func (a authInteractor) Refresh(ctx context.Context, tokenInput *input.TokenInput) (*output.TokenOutput, error) {
	authenticationErrors := app_error.NewAuthenticationErrors()

	if err := a.validate.StructCtx(ctx, tokenInput); err != nil {
		return nil, errors.WithStack(err)
	}

	standardClaim, err := a.tokenServiceFactory.GetTokenService(service.RefreshToken).DecryptToken(ctx, tokenInput.Token)
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
			detailErr := detail_error.NewTokenRevokedError()
			authenticationErrors.AddDetailError(detailErr)
			return nil, authenticationErrors
		}
		return nil, errors.WithStack(err)
	}

	if claim.GetJti().String() == jti.String() {
		detailErr := detail_error.NewTokenRevokedError()
		authenticationErrors.AddDetailError(detailErr)
		return nil, authenticationErrors
	}

	tokenOutput := new(output.TokenOutput)
	accessToken, err := a.tokenServiceFactory.GetTokenService(service.AccessToken).EncryptToken(ctx, standardClaim)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tokenOutput.AccessToken = accessToken

	return tokenOutput, nil
}

func (a authInteractor) Revoke(ctx context.Context, tokenInput *input.TokenInput) error {
	authenticationErrors := app_error.NewAuthenticationErrors()

	if err := a.validate.StructCtx(ctx, tokenInput); err != nil {
		return errors.WithStack(err)
	}

	standardClaim, err := a.tokenServiceFactory.GetTokenService(service.RefreshToken).DecryptToken(ctx, tokenInput.Token)
	if err != nil {
		return errors.WithStack(err)
	}

	userId, err := uuid.Parse(standardClaim.Subject)
	if err != nil {
		return errors.WithStack(err)
	}

	jti, err := uuid.Parse(standardClaim.Id)
	if err != nil {
		return errors.WithStack(err)
	}

	claim, err := a.claimService.GetById(ctx, userId)
	if err != nil {
		if app_error.IsRecordNotFoundError(errors.Cause(err)) {
			detailErr := detail_error.NewTokenRevokedError()
			authenticationErrors.AddDetailError(detailErr)
			return authenticationErrors
		}
		return errors.WithStack(err)
	}

	if claim.GetJti().String() == jti.String() {
		detailErr := detail_error.NewTokenRevokedError()
		authenticationErrors.AddDetailError(detailErr)
		return authenticationErrors
	}

	if err := a.claimService.Delete(ctx, claim); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
