package flags

import (
	"github.com/micro/cli/v2"
)

const refreshTokenSecret = "refresh_token_secret"

func NewRefreshTokenSecretFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    refreshTokenSecret,
		EnvVars: []string{"MICRO_REFRESH_TOKEN_SECRET"},
		Usage:   "Refresh token signing secret key for HMAC signing algorithm",
	}
}

func GetRefreshTokenSecret(ctx *cli.Context) string {
	return ctx.String(refreshTokenSecret)
}
