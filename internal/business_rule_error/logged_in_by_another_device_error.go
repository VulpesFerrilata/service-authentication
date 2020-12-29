package business_rule_error

import (
	"github.com/VulpesFerrilata/library/pkg/app_error"
	ut "github.com/go-playground/universal-translator"
)

func NewLoggedInByAnotherDeviceError() app_error.BusinessRuleError {
	return &LoggedInByAnotherDeviceError{}
}

type LoggedInByAnotherDeviceError struct{}

func (libade LoggedInByAnotherDeviceError) Error() string {
	return "your account has been logged in by another device"
}

func (libade LoggedInByAnotherDeviceError) Translate(trans ut.Translator) (string, error) {
	return trans.T("logged-in-by-another-device-error")
}
