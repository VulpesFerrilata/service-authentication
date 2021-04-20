package authentication_error

import (
	"fmt"
	"time"

	"github.com/VulpesFerrilata/library/pkg/app_error"
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

func (e expiredTokenError) Error() string {
	return fmt.Sprintf("token is expired by %v", e.delta)
}

func (e expiredTokenError) Translate(trans ut.Translator) (string, error) {
	detail, err := trans.T("expired-token-error", e.delta.String())
	return detail, errors.WithStack(err)
}

func (e expiredTokenError) ToAuthenticationErrors() app_error.AppError {
	return NewAuthenticationErrors(e)
}
