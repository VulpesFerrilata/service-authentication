package interactor

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/service"
	"github.com/VulpesFerrilata/auth/internal/usecase/adapter"
	"github.com/VulpesFerrilata/auth/internal/usecase/dto"
	"github.com/VulpesFerrilata/auth/internal/usecase/form"
	"github.com/VulpesFerrilata/grpc/protoc/user"
)

type AuthInteractor interface {
	Login(ctx context.Context, credentialRequest *user.CredentialRequest) (*dto.TokenDTO, error)
	Authenticate(ctx context.Context, tokenForm *form.TokenForm) (*dto.ClaimDTO, error)
	Refresh(ctx context.Context, tokenForm *form.TokenForm) (*dto.TokenDTO, error)
}

func NewAuthInteractor(authService service.AuthService, authAdapter adapter.AuthAdapter, userService user.UserService) AuthInteractor {
	return &authInteractor{
		authService: authService,
		authAdapter: authAdapter,
		userService: userService,
	}
}

type authInteractor struct {
	authService service.AuthService
	authAdapter adapter.AuthAdapter
	userService user.UserService
}

func (ai authInteractor) Login(ctx context.Context, credentialRequest *user.CredentialRequest) (*dto.TokenDTO, error) {
	userPb, err := ai.userService.GetUserByCredential(ctx, credentialRequest)
	if err != nil {
		return nil, err
	}

	token, err := ai.authAdapter.ParseUserPb(ctx, userPb)
	if err != nil {
		return nil, err
	}

	if err := ai.authService.CreateOrUpdate(ctx, token); err != nil {
		return nil, err
	}

	return ai.authAdapter.ResponseToken(ctx, token, false)
}

func (ai authInteractor) Authenticate(ctx context.Context, tokenForm *form.TokenForm) (*dto.ClaimDTO, error) {
	token, err := ai.authAdapter.ParseAccessToken(ctx, tokenForm)
	if err != nil {
		return nil, err
	}

	if err := ai.authService.ValidateAuthenticate(ctx, token); err != nil {
		return nil, err
	}

	return ai.authAdapter.ResponseClaim(ctx, token)
}

func (ai authInteractor) Refresh(ctx context.Context, tokenForm *form.TokenForm) (*dto.TokenDTO, error) {
	token, err := ai.authAdapter.ParseRefreshToken(ctx, tokenForm)
	if err != nil {
		return nil, err
	}

	if err := ai.authService.ValidateAuthenticate(ctx, token); err != nil {
		return nil, err
	}

	return ai.authAdapter.ResponseToken(ctx, token, true)
}
