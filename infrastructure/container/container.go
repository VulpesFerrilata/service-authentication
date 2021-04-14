package container

import (
	"reflect"
	"strings"

	"github.com/VulpesFerrilata/auth/internal/domain/repository"
	"github.com/VulpesFerrilata/auth/internal/domain/service"
	"github.com/VulpesFerrilata/auth/internal/pkg/micro/flags"
	"github.com/VulpesFerrilata/auth/internal/usecase/interactor"
	"github.com/VulpesFerrilata/grpc/gateway"
	"github.com/VulpesFerrilata/library/pkg/middleware"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/micro/cli/v2"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go.uber.org/dig"
)

func NewContainer(ctx *cli.Context) *dig.Container {
	sqlDialect := flags.GetSqlDialect(ctx)
	sqlDsn := flags.GetSqlDsn(ctx)
	translationFolderPath := flags.GetTranslationFolderPath(ctx)
	accessTokenAlg := flags.GetAccessTokenAlg(ctx)
	accessTokenSecret := flags.GetAccessTokenSecret(ctx)
	accessTokenDuration := flags.GetAccessTokenDuration(ctx)
	refreshTokenAlg := flags.GetRefreshTokenAlg(ctx)
	refreshTokenSecret := flags.GetRefreshTokenSecret(ctx)
	refreshTokenDuration := flags.GetRefreshTokenDuration(ctx)

	container := dig.New()

	container.Provide(func() (*gorm.DB, error) {
		var dialector gorm.Dialector
		switch strings.ToLower(sqlDialect) {
		case "mysql":
			dialector = mysql.Open(sqlDsn)
		case "sqlite":
			dialector = sqlite.Open(sqlDsn)
		default:
			err := errors.New("invalid sql dialect")
			return nil, errors.WithStack(err)
		}

		db, err := gorm.Open(dialector, &gorm.Config{})
		return db, errors.WithStack(err)
	})

	container.Provide(func() (*ut.UniversalTranslator, error) {
		en := en.New()
		utrans := ut.New(en, en)

		if err := utrans.Import(ut.FormatJSON, translationFolderPath); err != nil {
			return nil, errors.WithStack(err)
		}

		return utrans, errors.WithStack(utrans.VerifyTranslations())
	})

	container.Provide(func(utrans *ut.UniversalTranslator) (*validator.Validate, error) {
		v := validator.New()
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			jsonName := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if jsonName != "-" {
				return jsonName
			}
			return field.Name
		})

		en := en.New()
		trans, found := utrans.GetTranslator(en.Locale())
		if !found {
			err := errors.Errorf("translator not found: %v", en.Locale())
			return nil, errors.WithStack(err)
		}

		err := en_translations.RegisterDefaultTranslations(v, trans)
		return v, errors.WithStack(err)
	})

	//--Middleware
	container.Provide(middleware.NewRecoverMiddleware)
	container.Provide(middleware.NewTransactionMiddleware)
	container.Provide(middleware.NewTranslatorMiddleware)
	container.Provide(middleware.NewErrorHandlerMiddleware)

	//--Repositories
	container.Provide(repository.NewClaimRepository)
	//--Services
	container.Provide(service.NewUserService)
	container.Provide(service.NewClaimService)
	container.Provide(func() service.TokenServiceFactory {
		accessTokenService := service.NewTokenService(accessTokenAlg, accessTokenSecret, accessTokenDuration)
		refreshTokenService := service.NewTokenService(refreshTokenAlg, refreshTokenSecret, refreshTokenDuration)
		return service.NewTokenServiceFactory(accessTokenService, refreshTokenService)
	})
	//--Usecases
	container.Provide(interactor.NewAuthInteractor)
	//--Gateways
	container.Provide(func() gateway.UserGateway {
		return gateway.NewUserGateway()
	})

	return container
}
