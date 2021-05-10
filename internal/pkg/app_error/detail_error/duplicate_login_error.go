package detail_error

import (
	"github.com/VulpesFerrilata/library/pkg/app_error"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
)

func NewDuplicateLoginError() app_error.DetailError {
	return &duplicateLoginError{}
}

type duplicateLoginError struct{}

func (d duplicateLoginError) Error() string {
	return "user is currently logged in from another device"
}

func (d duplicateLoginError) Translate(trans ut.Translator) (string, error) {
	detail, err := trans.T("duplicate-login-error")
	return detail, errors.WithStack(err)
}
