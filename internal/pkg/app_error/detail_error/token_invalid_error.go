package detail_error

import (
	"github.com/VulpesFerrilata/library/pkg/app_error"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
)

func NewTokenInvalidError() app_error.DetailError {
	return &tokenInvalidError{}
}

type tokenInvalidError struct{}

func (t tokenInvalidError) Error() string {
	return "invalid token"
}

func (t tokenInvalidError) Translate(trans ut.Translator) (string, error) {
	detail, err := trans.T("invalid-token-error")
	return detail, errors.WithStack(err)
}
