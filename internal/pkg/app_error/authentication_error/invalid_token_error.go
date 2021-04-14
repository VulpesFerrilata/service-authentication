package authentication_error

import (
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
