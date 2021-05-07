package detail_error

import (
	"github.com/VulpesFerrilata/library/pkg/app_error"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
)

func NewTokenRevokedError() app_error.DetailError {
	return &tokenRevokedError{}
}

type tokenRevokedError struct{}

func (t tokenRevokedError) Error() string {
	return "token revoked"
}

func (t tokenRevokedError) Translate(trans ut.Translator) (string, error) {
	detail, err := trans.T("token-revoked-error")
	return detail, errors.WithStack(err)
}
