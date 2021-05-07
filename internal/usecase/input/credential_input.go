package input

type CredentialInput struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}
