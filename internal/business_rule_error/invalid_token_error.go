package business_rule_error

import (
	"github.com/VulpesFerrilata/library/pkg/app_error"
	ut "github.com/go-playground/universal-translator"
)

func NewInvalidTokenError() app_error.BusinessRuleError {
	return &invalidTokenError{}
}

type invalidTokenError struct{}

func (ite invalidTokenError) Error() string {
	return "invalid token"
}

func (ite invalidTokenError) Translate(trans ut.Translator) (string, error) {
	return trans.T("invalid-token-error")
}
