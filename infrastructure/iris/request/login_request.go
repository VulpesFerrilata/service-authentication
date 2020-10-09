package request

import (
	"github.com/VulpesFerrilata/grpc/protoc/user"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (lr LoginRequest) ToCredentialRequestPb() *user.CredentialRequest {
	credentialRequestPb := new(user.CredentialRequest)
	credentialRequestPb.Username = lr.Username
	credentialRequestPb.Password = lr.Password
	return credentialRequestPb
}
