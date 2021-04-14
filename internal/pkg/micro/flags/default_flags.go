package flags

import (
	"time"

	"github.com/micro/cli/v2"
)

const (
	sqlDialect            = "sql_dialect"
	sqlDsn                = "sql_dsn"
	translationFolderPath = "translation_folder_path"
	accessTokenAlg        = "access_token_alg"
	accessTokenSecret     = "access_token_secret"
	accessTokenDuration   = "access_token_duration"
	refreshTokenAlg       = "refresh_token_alg"
	refreshTokenSecret    = "refresh_token_secret"
	refreshTokenDuration  = "refresh_token_duration"
)

var DefaultFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    sqlDialect,
		EnvVars: []string{"MICRO_SQL_DIALECT"},
		Usage:   "Sql dialect for storing data, currently support sqlite and mysql",
	},
	&cli.StringFlag{
		Name:    sqlDsn,
		EnvVars: []string{"MICRO_SQL_DSN"},
		Usage:   "Connection string which used for sql dialect to initialize data source",
	},
	&cli.StringFlag{
		Name:    translationFolderPath,
		EnvVars: []string{"MICRO_TRANSLATION_FOLDER_PATH"},
		Usage:   "Folder path contains translation information which will be imported by universal translator",
	},
	&cli.StringFlag{
		Name:    sqlDialect,
		EnvVars: []string{"MICRO_SQL_DIALECT"},
		Usage:   "Sql dialect for storing data, currently support sqlite and mysql",
	},
}

func GetSqlDialect(ctx *cli.Context) string {
	return ctx.String(sqlDialect)
}

func GetSqlDsn(ctx *cli.Context) string {
	return ctx.String(sqlDsn)
}

func GetTranslationFolderPath(ctx *cli.Context) string {
	return ctx.String(translationFolderPath)
}

func GetAccessTokenAlg(ctx *cli.Context) string {
	return ctx.String(accessTokenAlg)
}

func GetAccessTokenSecret(ctx *cli.Context) string {
	return ctx.String(accessTokenSecret)
}

func GetAccessTokenDuration(ctx *cli.Context) time.Duration {
	return ctx.Duration(accessTokenDuration)
}

func GetRefreshTokenAlg(ctx *cli.Context) string {
	return ctx.String(refreshTokenAlg)
}

func GetRefreshTokenSecret(ctx *cli.Context) string {
	return ctx.String(refreshTokenSecret)
}

func GetRefreshTokenDuration(ctx *cli.Context) time.Duration {
	return ctx.Duration(refreshTokenDuration)
}
