package service

import (
	"context"
	"time"

	"github.com/VulpesFerrilata/auth/internal/pkg/app_error/authentication_error"
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
	parser := &jwt.Parser{
		SkipClaimsValidation: true,
		ValidMethods:         []string{t.alg},
	}

	standardClaim := new(jwt.StandardClaims)
	_, err := parser.ParseWithClaims(token, standardClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secretKey), nil
	})
	if err != nil {
		return nil, authentication_error.NewInvalidTokenError()
	}

	if err := t.validate(ctx, standardClaim); err != nil {
		return nil, err
	}

	return standardClaim, nil
}

func (t tokenService) validate(ctx context.Context, standardClaim *jwt.StandardClaims) error {
	now := time.Now().Unix()
	if !standardClaim.VerifyExpiresAt(now, true) {
		delta := time.Unix(now, 0).Sub(time.Unix(standardClaim.ExpiresAt, 0))
		return authentication_error.NewExpiredTokenError(delta)
	}

	return nil
}
