package authentication_error

import (
	"strings"

	"github.com/VulpesFerrilata/library/pkg/app_error"
	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewAuthenticationErrors(authenticationErrs ...AuthenticationError) app_error.AppError {
	return authenticationErrors(authenticationErrs)
}

type authenticationErrors []AuthenticationError

func (a authenticationErrors) Error() string {
	builder := new(strings.Builder)

	builder.WriteString("authentication failed")
	for _, authenticationErr := range a {
		builder.WriteString("\n")
		builder.WriteString(authenticationErr.Error())
	}

	return builder.String()
}

func (a authenticationErrors) Problem(trans ut.Translator) (iris.Problem, error) {
	problem := iris.NewProblem()
	problem.Type("about:blank")

	problem.Status(iris.StatusUnprocessableEntity)
	detail, err := trans.T("authentication-error")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	problem.Detail(detail)

	var errs []string
	for _, authenticationErr := range a {
		authenticationErrTrans, err := authenticationErr.Translate(trans)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		errs = append(errs, authenticationErrTrans)
	}
	problem.Key("errors", errs)

	return problem, nil
}

func (a authenticationErrors) Status(trans ut.Translator) (*status.Status, error) {
	detail, err := trans.T("authentication-error")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	stt := status.New(codes.Unauthenticated, detail)

	preconditionFailure := &errdetails.PreconditionFailure{}
	for _, authenticationErr := range a {
		authenticationErrTrans, err := authenticationErr.Translate(trans)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		violation := &errdetails.PreconditionFailure_Violation{
			Type:        "AUTHENTICATION",
			Description: authenticationErrTrans,
		}

		preconditionFailure.Violations = append(preconditionFailure.Violations, violation)
	}

	stt, err = stt.WithDetails(preconditionFailure)
	return stt, errors.WithStack(err)
}

func (a authenticationErrors) Message(trans ut.Translator) (string, error) {
	builder := new(strings.Builder)

	detail, err := trans.T("authentication-error")
	if err != nil {
		return "", errors.WithStack(err)
	}
	builder.WriteString(detail)
	for _, authenticationErr := range a {
		builder.WriteString("\n")
		authenticationErrTrans, err := authenticationErr.Translate(trans)
		if err != nil {
			return "", errors.WithStack(err)
		}
		builder.WriteString(authenticationErrTrans)
	}

	return builder.String(), nil
}
