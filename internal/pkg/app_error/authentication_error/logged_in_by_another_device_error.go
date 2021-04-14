package authentication_error

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
)

func NewLoggedInByAnotherDeviceError() AuthenticationError {
	return &loggedInByAnotherDeviceError{}
}

type loggedInByAnotherDeviceError struct{}

func (l loggedInByAnotherDeviceError) Error() string {
	return "your account has been logged in by another device"
}

func (l loggedInByAnotherDeviceError) Translate(trans ut.Translator) (string, error) {
	detail, err := trans.T("logged-in-by-another-device-error")
	return detail, errors.WithStack(err)
}
