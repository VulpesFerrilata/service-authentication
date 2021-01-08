package authentication_error

import (
	"fmt"
	"time"

	"github.com/VulpesFerrilata/library/pkg/app_error"
	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewExpiredTokenError(delta time.Duration) app_error.AppError {
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

func (ete expiredTokenError) Problem(trans ut.Translator) (iris.Problem, error) {
	problem := iris.NewProblem()
	problem.Type("about:blank")
	problem.Status(iris.StatusUnauthorized)

	detail, err := trans.T("expired-token-error", ete.delta.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, "expired-token-error")
	}
	problem.Detail(detail)

	return problem, nil
}

func (ete expiredTokenError) Status(trans ut.Translator) (*status.Status, error) {
	detail, err := trans.T("expired-token-error", ete.delta.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, "expired-token-error")
	}
	stt := status.New(codes.Unauthenticated, detail)
	return stt, nil
}

func (ete expiredTokenError) Message(trans ut.Translator) (string, error) {
	msg, err := trans.T("expired-token-error", ete.delta.String())
	return msg, fmt.Errorf("%w: %s", err, "expired-token-error")
}
