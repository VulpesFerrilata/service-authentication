package request

type UserCredentialRequest struct {
	ID       string `json:"id" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
