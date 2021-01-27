package authentication_error

import (
	"fmt"

	"github.com/VulpesFerrilata/library/pkg/app_error"
	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewLoggedInByAnotherDeviceError() app_error.AppError {
	return &loggedInByAnotherDeviceError{}
}

type loggedInByAnotherDeviceError struct{}

func (libade loggedInByAnotherDeviceError) Error() string {
	return "your account has been logged in by another device"
}

func (libade loggedInByAnotherDeviceError) Problem(trans ut.Translator) (iris.Problem, error) {
	problem := iris.NewProblem()
	problem.Type("about:blank")
	problem.Status(iris.StatusUnauthorized)

	detail, err := trans.T("logged-in-by-another-device-error")
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, "expired-token-error")
	}
	problem.Detail(detail)

	return problem, nil
}

func (libade loggedInByAnotherDeviceError) Status(trans ut.Translator) (*status.Status, error) {
	detail, err := trans.T("logged-in-by-another-device-error")
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, "expired-token-error")
	}
	stt := status.New(codes.Unauthenticated, detail)
	return stt, nil
}

func (libade loggedInByAnotherDeviceError) Message(trans ut.Translator) (string, error) {
	msg, err := trans.T("logged-in-by-another-device-error")
	return msg, fmt.Errorf("%w: %s", err, "expired-token-error")
}
