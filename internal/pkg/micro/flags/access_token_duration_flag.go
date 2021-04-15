package flags

import (
	"time"

	"github.com/micro/cli/v2"
)

const accessTokenDuration = "access_token_duration"

func NewAccessTokenDurationFlag() cli.Flag {
	return &cli.DurationFlag{
		Name:    accessTokenDuration,
		EnvVars: []string{"MICRO_ACCESS_TOKEN_DURATION"},
		Usage:   "Access token signing secret key for HMAC signing algorithm",
	}
}

func GetAccessTokenDuration(ctx *cli.Context) time.Duration {
	return ctx.Duration(accessTokenDuration)
}
