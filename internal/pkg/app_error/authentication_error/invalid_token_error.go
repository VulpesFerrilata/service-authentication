package authentication_error

import (
	"fmt"

	"github.com/VulpesFerrilata/library/pkg/app_error"
	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewInvalidTokenError() app_error.AppError {
	return &invalidTokenError{}
}

type invalidTokenError struct{}

func (ite invalidTokenError) Error() string {
	return "invalid token"
}

func (ite invalidTokenError) Problem(trans ut.Translator) (iris.Problem, error) {
	problem := iris.NewProblem()
	problem.Type("about:blank")
	problem.Status(iris.StatusUnauthorized)

	detail, err := trans.T("invalid-token-error")
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, "expired-token-error")
	}
	problem.Detail(detail)

	return problem, nil
}

func (ite invalidTokenError) Status(trans ut.Translator) (*status.Status, error) {
	detail, err := trans.T("invalid-token-error")
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, "expired-token-error")
	}
	stt := status.New(codes.Unauthenticated, detail)
	return stt, nil
}

func (ite invalidTokenError) Message(trans ut.Translator) (string, error) {
	msg, err := trans.T("invalid-token-error")
	return msg, fmt.Errorf("%w: %s", err, "expired-token-error")
}
