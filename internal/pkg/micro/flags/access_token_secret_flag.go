package flags

import (
	"github.com/micro/cli/v2"
)

const accessTokenSecret = "access_token_secret"

func NewAccessTokenSecretFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    accessTokenSecret,
		EnvVars: []string{"MICRO_ACCESS_TOKEN_SECRET"},
		Usage:   "Access token signing secret key for HMAC signing algorithm",
	}
}

func GetAccessTokenSecret(ctx *cli.Context) string {
	return ctx.String(accessTokenSecret)
}
