package detail_error

import (
	"github.com/VulpesFerrilata/library/pkg/app_error"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
)

func NewInvalidPasswordError() app_error.DetailError {
	return &invalidPasswordError{}
}

type invalidPasswordError struct{}

func (i invalidPasswordError) Error() string {
	return "password is invalid"
}

func (i invalidPasswordError) Translate(trans ut.Translator) (string, error) {
	detail, err := trans.T("invalid-password-error")
	return detail, errors.WithStack(err)
}
