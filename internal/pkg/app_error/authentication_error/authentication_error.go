package authentication_error

import (
	"github.com/VulpesFerrilata/library/pkg/app_error"
	ut "github.com/go-playground/universal-translator"
)

type AuthenticationError interface {
	error
	Translate(trans ut.Translator) (string, error)
	ToAuthenticationErrors() app_error.AppError
}
