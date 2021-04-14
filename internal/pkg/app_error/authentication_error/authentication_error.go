package authentication_error

import ut "github.com/go-playground/universal-translator"

type AuthenticationError interface {
	error
	Translate(trans ut.Translator) (string, error)
}
