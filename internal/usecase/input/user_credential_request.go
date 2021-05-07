package input

type UserCredentialInput struct {
	ID       string `validate:"required"`
	Username string `validate:"required"`
	Password string `validate:"required"`
}
