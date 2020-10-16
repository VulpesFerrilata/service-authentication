package request

type TokenRequest struct {
	Token string `json:"token" validate:"required"`
}