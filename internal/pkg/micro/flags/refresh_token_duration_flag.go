package flags

import (
	"time"

	"github.com/micro/cli/v2"
)

const refreshTokenDuration = "refresh_token_duration"

func NewRefreshTokenDurationFlag() cli.Flag {
	return &cli.DurationFlag{
		Name:    refreshTokenDuration,
		EnvVars: []string{"MICRO_REFRESH_TOKEN_DURATION"},
		Usage:   "Refresh token signing secret key for HMAC signing algorithm",
	}
}

func GetRefreshTokenDuration(ctx *cli.Context) time.Duration {
	return ctx.Duration(refreshTokenDuration)
}
