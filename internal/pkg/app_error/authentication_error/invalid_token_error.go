package authentication_error

import (
	"github.com/VulpesFerrilata/library/pkg/app_error"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
)

func NewInvalidTokenError() AuthenticationError {
	return &invalidTokenError{}
}

type invalidTokenError struct{}

func (i invalidTokenError) Error() string {
	return "invalid token"
}

func (i invalidTokenError) Translate(trans ut.Translator) (string, error) {
	detail, err := trans.T("invalid-token-error")
	return detail, errors.WithStack(err)
}

func (i invalidTokenError) ToAuthenticationErrors() app_error.AppError {
	return NewAuthenticationErrors(i)
}
