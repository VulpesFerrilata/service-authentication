package flags

import (
	"github.com/VulpesFerrilata/library/pkg/micro/flags/generic"
	"github.com/micro/cli/v2"
)

const refreshTokenAlg = "refresh_token_alg"

func NewRefreshTokenAlgFlag() cli.Flag {
	return &cli.GenericFlag{
		Name:    refreshTokenAlg,
		Value:   generic.NewStringGeneric("HS256", "HS384", "HS512"),
		EnvVars: []string{"MICRO_REFRESH_TOKEN_ALG"},
		Usage:   "Refresh token signing algorithm (HMAC), currently support HS256, HS384 and HS512",
	}
}

func GetRefreshTokenAlg(ctx *cli.Context) string {
	return ctx.String(refreshTokenAlg)
}
