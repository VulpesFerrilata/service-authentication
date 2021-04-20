package flags

import (
	"github.com/VulpesFerrilata/library/pkg/micro/flags/generic"
	"github.com/micro/cli/v2"
)

const accessTokenAlg = "access_token_alg"

func NewAccessTokenAlgFlag() cli.Flag {
	return &cli.GenericFlag{
		Name:    accessTokenAlg,
		Value:   generic.NewStringGeneric("HS256", "HS384", "HS512"),
		EnvVars: []string{"MICRO_ACCESS_TOKEN_ALG"},
		Usage:   "Access token signing algorithm (HMAC), currently support HS256, HS384 and HS512",
	}
}

func GetAccessTokenAlg(ctx *cli.Context) string {
	return ctx.Generic(accessTokenAlg).(cli.Generic).String()
}
