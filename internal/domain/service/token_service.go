package service

import (
	"context"
	"time"

	"github.com/VulpesFerrilata/auth/internal/pkg/app_error/detail_error"
	"github.com/VulpesFerrilata/library/pkg/app_error"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type TokenService interface {
	EncryptToken(ctx context.Context, standardClaim *jwt.StandardClaims) (string, error)
	DecryptToken(ctx context.Context, token string) (*jwt.StandardClaims, error)
}

func NewTokenService(alg string, secretKey string, duration time.Duration) TokenService {
	return &tokenService{
		alg:       alg,
		secretKey: secretKey,
		duration:  duration,
	}
}

type tokenService struct {
	alg       string
	secretKey string
	duration  time.Duration
}

func (t tokenService) EncryptToken(ctx context.Context, standardClaim *jwt.StandardClaims) (string, error) {
	standardClaim.IssuedAt = time.Now().Unix()
	standardClaim.ExpiresAt = time.Now().Add(t.duration).Unix()
	token, err := jwt.NewWithClaims(jwt.GetSigningMethod(t.alg), standardClaim).SignedString([]byte(t.secretKey))
	return token, errors.WithStack(err)
}

func (t tokenService) DecryptToken(ctx context.Context, token string) (*jwt.StandardClaims, error) {
	authenticationErrs := app_error.NewAuthenticationErrors()

	parser := &jwt.Parser{
		SkipClaimsValidation: true,
		ValidMethods:         []string{t.alg},
	}

	standardClaim := new(jwt.StandardClaims)
	if _, err := parser.ParseWithClaims(token, standardClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secretKey), nil
	}); err != nil {
		detailErr := detail_error.NewTokenInvalidError()
		authenticationErrs.AddDetailError(detailErr)
		return nil, authenticationErrs
	}

	if err := t.validate(ctx, standardClaim); err != nil {
		return nil, errors.WithStack(err)
	}

	return standardClaim, nil
}

func (t tokenService) validate(ctx context.Context, standardClaim *jwt.StandardClaims) error {
	authenticationErrs := app_error.NewAuthenticationErrors()

	now := time.Now().Unix()
	if !standardClaim.VerifyExpiresAt(now, true) {
		delta := time.Unix(now, 0).Sub(time.Unix(standardClaim.ExpiresAt, 0))
		detailErr := detail_error.NewTokenExpiredError(delta)
		authenticationErrs.AddDetailError(detailErr)
		return authenticationErrs
	}

	return nil
}
