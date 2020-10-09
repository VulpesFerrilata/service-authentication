package adapter

import (
	"context"
	"time"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/auth/internal/usecase/dto"
	"github.com/VulpesFerrilata/auth/internal/usecase/form"
	"github.com/VulpesFerrilata/grpc/protoc/user"
	"github.com/VulpesFerrilata/library/config"
	"github.com/VulpesFerrilata/library/pkg/validator"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/iris-contrib/go.uuid"
)

type AuthAdapter interface {
	ParseUserPb(ctx context.Context, userPb *user.UserResponse) (*model.Token, error)
	ParseAccessToken(ctx context.Context, tokenForm *form.TokenForm) (*model.Token, error)
	ParseRefreshToken(ctx context.Context, tokenForm *form.TokenForm) (*model.Token, error)
	ResponseToken(ctx context.Context, token *model.Token, accessTokenOnly bool) (*dto.TokenDTO, error)
	ResponseClaim(ctx context.Context, token *model.Token) (*dto.ClaimDTO, error)
}

func NewAuthAdapter(jwtCfg *config.JwtConfig, validate validator.Validate) AuthAdapter {
	return &authAdapter{
		jwtCfg:   jwtCfg,
		validate: validate,
	}
}

type authAdapter struct {
	jwtCfg   *config.JwtConfig
	validate validator.Validate
}

func (ap authAdapter) ParseUserPb(ctx context.Context, userPb *user.UserResponse) (*model.Token, error) {
	token := new(model.Token)
	token.UserID = uint(userPb.ID)

	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	token.Jti = uuid.String()

	return token, nil
}

func (ap authAdapter) ParseAccessToken(ctx context.Context, tokenForm *form.TokenForm) (*model.Token, error) {
	if err := ap.validate.Struct(ctx, tokenForm); err != nil {
		return nil, err
	}
	return tokenForm.ToToken(ap.jwtCfg.AccessTokenSettings)
}

func (ap authAdapter) ParseRefreshToken(ctx context.Context, tokenForm *form.TokenForm) (*model.Token, error) {
	if err := ap.validate.Struct(ctx, tokenForm); err != nil {
		return nil, err
	}
	return tokenForm.ToToken(ap.jwtCfg.RefreshTokenSettings)
}

func (ap authAdapter) ResponseToken(ctx context.Context, token *model.Token, accessTokenOnly bool) (*dto.TokenDTO, error) {
	tokenDTO := new(dto.TokenDTO)

	accessToken, err := ap.createToken(token, ap.jwtCfg.AccessTokenSettings)
	if err != nil {
		return nil, err
	}
	tokenDTO.AccessToken = accessToken

	if !accessTokenOnly {
		refreshToken, err := ap.createToken(token, ap.jwtCfg.RefreshTokenSettings)
		if err != nil {
			return nil, err
		}
		tokenDTO.RefreshToken = refreshToken
	}

	return tokenDTO, nil
}

func (ap authAdapter) ResponseClaim(ctx context.Context, token *model.Token) (*dto.ClaimDTO, error) {
	ClaimDTO := new(dto.ClaimDTO)
	ClaimDTO.UserID = int(token.UserID)
	return ClaimDTO, nil
}

func (ap authAdapter) createToken(token *model.Token, tokenSettings config.TokenSettings) (string, error) {
	claim := new(jwt.StandardClaims)
	claim.Id = token.Jti
	claim.Subject = string(token.UserID)
	claim.IssuedAt = time.Now().Unix()
	claim.ExpiresAt = time.Now().Add(tokenSettings.Duration).Unix()

	return jwt.NewWithClaims(jwt.GetSigningMethod(tokenSettings.Alg), claim).SignedString([]byte(tokenSettings.SecretKey))
}
