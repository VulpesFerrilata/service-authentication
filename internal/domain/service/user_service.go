package service

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/grpc/protoc/user"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UserService interface {
	GetByCredential(ctx context.Context, username string, password string) (*model.User, error)
}

func NewUserService(userGateway user.UserService) UserService {
	return &userService{
		userGateway: userGateway,
	}
}

type userService struct {
	userGateway user.UserService
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

	userId, err := uuid.Parse(userResponsePb.GetID())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	user := model.NewUser(
		userId,
		userResponsePb.GetUsername(),
		userResponsePb.GetDisplayName(),
		userResponsePb.GetEmail(),
	)

	return user, nil
}
