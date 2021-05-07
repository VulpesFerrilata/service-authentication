package detail_error

import (
	"fmt"
	"time"

	"github.com/VulpesFerrilata/library/pkg/app_error"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
)

func NewTokenExpiredError(delta time.Duration) app_error.DetailError {
	return &tokenExpiredError{
		delta: delta,
	}
}

type tokenExpiredError struct {
	delta time.Duration
}

func (t tokenExpiredError) Error() string {
	return fmt.Sprintf("token is expired by %v", t.delta)
}

func (t tokenExpiredError) Translate(trans ut.Translator) (string, error) {
	detail, err := trans.T("expired-token-error", t.delta.String())
	return detail, errors.WithStack(err)
}
