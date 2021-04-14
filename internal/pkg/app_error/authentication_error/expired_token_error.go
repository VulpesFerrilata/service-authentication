package authentication_error

import (
	"fmt"
	"time"

	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
)

func NewExpiredTokenError(delta time.Duration) AuthenticationError {
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
	detail, err := trans.T("expired-token-error", ete.delta.String())
	return detail, errors.WithStack(err)
}
