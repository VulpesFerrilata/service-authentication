package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
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
	EncryptToken(ctx context.Context, tokenType TokenType, claim *model.Claim) (string, error)
	DecryptToken(ctx context.Context, tokenType TokenType, token string) (*model.Claim, error)
}

func NewTokenService(jwtCfg *config.JwtConfig) TokenService {
	return &tokenService{
		jwtCfg: jwtCfg,
	}
}

type tokenService struct {
	jwtCfg *config.JwtConfig
}

func (ts tokenService) EncryptToken(ctx context.Context, tokenType TokenType, claim *model.Claim) (string, error) {
	switch tokenType {
	case AccessToken:
		return ts.encryptToken(ctx, claim, ts.jwtCfg.AccessTokenSettings)
	case RefreshToken:
		return ts.encryptToken(ctx, claim, ts.jwtCfg.RefreshTokenSettings)
	default:
		return "", InvalidTokenTypeErr
	}
}

func (ts tokenService) DecryptToken(ctx context.Context, tokenType TokenType, token string) (*model.Claim, error) {
	switch tokenType {
	case AccessToken:
		return ts.decryptToken(ctx, token, ts.jwtCfg.AccessTokenSettings)
	case RefreshToken:
		return ts.decryptToken(ctx, token, ts.jwtCfg.RefreshTokenSettings)
	default:
		return nil, InvalidTokenTypeErr
	}
}

func (ts tokenService) encryptToken(ctx context.Context, claim *model.Claim, tokenSettings config.TokenSettings) (string, error) {
	standardClaim := new(jwt.StandardClaims)
	standardClaim.Id = claim.GetJti()
	standardClaim.Subject = fmt.Sprint(claim.GetUserId())
	standardClaim.IssuedAt = time.Now().Unix()
	standardClaim.ExpiresAt = time.Now().Add(tokenSettings.Duration).Unix()

	return jwt.NewWithClaims(jwt.GetSigningMethod(tokenSettings.Alg), standardClaim).SignedString([]byte(tokenSettings.SecretKey))
}

func (ts tokenService) decryptToken(ctx context.Context, token string, tokenSettings config.TokenSettings) (*model.Claim, error) {
	parser := &jwt.Parser{
		SkipClaimsValidation: true,
	}

	standardClaim := new(jwt.StandardClaims)
	_, err := parser.ParseWithClaims(token, standardClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSettings.SecretKey), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "service.TokenService.decryptToken")
	}

	//validate

	userId, err := strconv.ParseUint(standardClaim.Subject, 10, 32)
	if err != nil {
		return nil, errors.Wrap(err, "service.TokenService.decryptToken")
	}
	claim := model.NewClaimWithJti(uint(userId), standardClaim.Id)

	return claim, nil
}
