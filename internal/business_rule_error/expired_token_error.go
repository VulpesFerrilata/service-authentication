package business_rule_error

import (
	"fmt"
	"time"

	"github.com/VulpesFerrilata/library/pkg/app_error"
	ut "github.com/go-playground/universal-translator"
)

func NewExpiredTokenError(delta time.Duration) app_error.BusinessRuleError {
	return &expiredTokenError{
		delta: delta,
	}
}

type expiredTokenError struct {
	delta time.Duration
}

func (ete expiredTokenError) Error() string {
	return fmt.Sprintf("token is expired by %v", ete.delta)
}

func (ete expiredTokenError) Translate(trans ut.Translator) (string, error) {
	return trans.T("expired-token-error", ete.delta.String())
}
