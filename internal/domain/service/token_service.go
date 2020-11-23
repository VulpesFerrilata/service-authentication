package service

import (
	"context"
	"strconv"
	"time"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/library/config"
	"github.com/dgrijalva/jwt-go"
)

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
)

type TokenService interface {
	EncryptAccessToken(ctx context.Context, claim *model.Claim) (string, error)
	EncryptRefreshToken(ctx context.Context, claim *model.Claim) (string, error)
	DecryptAccessToken(ctx context.Context, token string) (*model.Claim, error)
	DecryptRefreshToken(ctx context.Context, token string) (*model.Claim, error)
}

func NewTokenService(jwtCfg *config.JwtConfig) TokenService {
	return &tokenService{
		jwtCfg: jwtCfg,
	}
}

type tokenService struct {
	jwtCfg *config.JwtConfig
}

func (ts tokenService) EncryptAccessToken(ctx context.Context, claim *model.Claim) (string, error) {
	return ts.encryptToken(ctx, claim, ts.jwtCfg.AccessTokenSettings)
}

func (ts tokenService) EncryptRefreshToken(ctx context.Context, claim *model.Claim) (string, error) {
	return ts.encryptToken(ctx, claim, ts.jwtCfg.RefreshTokenSettings)
}

func (ts tokenService) DecryptAccessToken(ctx context.Context, token string) (*model.Claim, error) {
	return ts.decryptToken(ctx, token, ts.jwtCfg.AccessTokenSettings)
}
func (ts tokenService) DecryptRefreshToken(ctx context.Context, token string) (*model.Claim, error) {
	return ts.decryptToken(ctx, token, ts.jwtCfg.RefreshTokenSettings)
}

func (ts tokenService) encryptToken(ctx context.Context, claim *model.Claim, tokenSettings config.TokenSettings) (string, error) {
	standardClaim := new(jwt.StandardClaims)
	standardClaim.Id = claim.Jti
	standardClaim.Subject = string(claim.UserID)
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
		return nil, err
	}

	//validate

	claim := new(model.Claim)
	claim.Jti = standardClaim.Id

	userId, err := strconv.ParseUint(standardClaim.Subject, 10, 32)
	if err != nil {
		return nil, err
	}
	claim.UserID = uint(userId)

	return claim, nil
}
