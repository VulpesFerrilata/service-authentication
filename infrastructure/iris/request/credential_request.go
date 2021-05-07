package request

type CredentialRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
