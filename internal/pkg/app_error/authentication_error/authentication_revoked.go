package authentication_error

import (
	"github.com/VulpesFerrilata/library/pkg/app_error"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
)

func NewAuthenticationRevokedError() AuthenticationError {
	return &authenticationRevokedError{}
}

type authenticationRevokedError struct{}

func (a authenticationRevokedError) Error() string {
	return "authentication revoked"
}

func (a authenticationRevokedError) Translate(trans ut.Translator) (string, error) {
	detail, err := trans.T("authentication-revoked-error")
	return detail, errors.WithStack(err)
}

func (a authenticationRevokedError) ToAuthenticationErrors() app_error.AppError {
	return NewAuthenticationErrors(a)
}
