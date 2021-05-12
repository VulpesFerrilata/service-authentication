package input

type UserCredentialInput struct {
	UserID   string `validate:"required"`
	Password string `validate:"required"`
}
