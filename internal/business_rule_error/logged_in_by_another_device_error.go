package business_rule_error

import (
	"github.com/VulpesFerrilata/library/pkg/app_error"
	ut "github.com/go-playground/universal-translator"
)

func NewLoggedInByAnotherDeviceError() app_error.BusinessRuleError {
	return &loggedInByAnotherDeviceError{}
}

type loggedInByAnotherDeviceError struct{}

func (libade loggedInByAnotherDeviceError) Error() string {
	return "your account has been logged in by another device"
}

func (libade loggedInByAnotherDeviceError) Translate(trans ut.Translator) (string, error) {
	return trans.T("logged-in-by-another-device-error")
}
