package request

import "github.com/VulpesFerrilata/grpc/protoc/user"

type CredentialRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (cr CredentialRequest) ToCredentialRequestPb() *user.CredentialRequest {
	credentialRequestPb := new(user.CredentialRequest)
	credentialRequestPb.Username = cr.Username
	credentialRequestPb.Password = cr.Password
	return credentialRequestPb
}
