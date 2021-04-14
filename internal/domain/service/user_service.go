package service

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/auth/internal/mapper"
	"github.com/VulpesFerrilata/grpc/gateway"
	"github.com/VulpesFerrilata/grpc/protoc/user"
	"github.com/pkg/errors"
)

type UserService interface {
	GetByCredential(ctx context.Context, username string, password string) (*model.User, error)
}

func NewUserService(userGateway gateway.UserGateway) UserService {
	return &userService{
		userGateway: userGateway,
	}
}

type userService struct {
	userGateway gateway.UserGateway
}

func (u userService) GetByCredential(ctx context.Context, username string, password string) (*model.User, error) {
	credentialRequestPb := user.CredentialRequest{
		Username: username,
		Password: password,
	}

	userResponsePb, err := u.userGateway.GetUserByCredential(ctx, &credentialRequestPb)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	user, err := mapper.NewUserResponsePbMapper(userResponsePb).ToUserModel()
	return user, errors.WithStack(err)
}
