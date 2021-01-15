package service

import (
	"context"
	"fmt"
	"time"

	"github.com/VulpesFerrilata/auth/internal/app_error/authentication_error"
	"github.com/VulpesFerrilata/auth/internal/domain/datamodel"
	"github.com/VulpesFerrilata/library/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
)

var InvalidTokenTypeErr = errors.New("invalid token type")

type TokenService interface {
	EncryptToken(ctx context.Context, tokenType TokenType, claim *datamodel.Claim) (string, error)
	DecryptToken(ctx context.Context, tokenType TokenType, token string) (*datamodel.Claim, error)
}

func NewTokenService(jwtCfg *config.JwtConfig) TokenService {
	return &tokenService{
		jwtCfg: jwtCfg,
	}
}

type tokenService struct {
	jwtCfg *config.JwtConfig
}

func (ts tokenService) EncryptToken(ctx context.Context, tokenType TokenType, claim *datamodel.Claim) (string, error) {
	switch tokenType {
	case AccessToken:
		return ts.encryptToken(ctx, claim, ts.jwtCfg.AccessTokenSettings)
	case RefreshToken:
		return ts.encryptToken(ctx, claim, ts.jwtCfg.RefreshTokenSettings)
	default:
		return "", InvalidTokenTypeErr
	}
}

func (ts tokenService) DecryptToken(ctx context.Context, tokenType TokenType, token string) (*datamodel.Claim, error) {
	switch tokenType {
	case AccessToken:
		return ts.decryptToken(ctx, token, ts.jwtCfg.AccessTokenSettings)
	case RefreshToken:
		return ts.decryptToken(ctx, token, ts.jwtCfg.RefreshTokenSettings)
	default:
		return nil, InvalidTokenTypeErr
	}
}

func (ts tokenService) encryptToken(ctx context.Context, claim *datamodel.Claim, tokenSettings config.TokenSettings) (string, error) {
	standardClaim := new(jwt.StandardClaims)
	standardClaim.Id = claim.GetJti().String()
	standardClaim.Subject = fmt.Sprint(claim.GetUserId())
	standardClaim.IssuedAt = time.Now().Unix()
	standardClaim.ExpiresAt = time.Now().Add(tokenSettings.Duration).Unix()

	return jwt.NewWithClaims(jwt.GetSigningMethod(tokenSettings.Alg), standardClaim).SignedString([]byte(tokenSettings.SecretKey))
}

func (ts tokenService) decryptToken(ctx context.Context, token string, tokenSettings config.TokenSettings) (*datamodel.Claim, error) {
	parser := &jwt.Parser{
		SkipClaimsValidation: true,
		ValidMethods:         []string{tokenSettings.Alg},
	}

	standardClaim := new(jwt.StandardClaims)
	_, err := parser.ParseWithClaims(token, standardClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSettings.SecretKey), nil
	})
	if err != nil {
		return nil, authentication_error.NewInvalidTokenError()
	}

	//validate
	now := time.Now().Unix()
	if !standardClaim.VerifyExpiresAt(now, true) {
		delta := time.Unix(now, 0).Sub(time.Unix(standardClaim.ExpiresAt, 0))
		return nil, authentication_error.NewExpiredTokenError(delta)
	}

	return datamodel.NewClaimFromStandardClaim(standardClaim)
}
