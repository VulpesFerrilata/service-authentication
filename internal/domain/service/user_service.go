package service

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/grpc/protoc/user"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UserService interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}

func NewUser(userGateway user.UserService) UserService {
	return &userService{
		userGateway: userGateway,
	}
}

type userService struct {
	userGateway user.UserService
}

func (u userService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	usernameRequestPb := &user.UsernameRequest{
		Username: username,
	}

	userResponsePb, err := u.userGateway.GetUserByUsername(ctx, usernameRequestPb)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	id, err := uuid.Parse(userResponsePb.GetID())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	user := model.NewUser(
		id,
		userResponsePb.GetDisplayName(),
		userResponsePb.GetEmail(),
	)

	return user, nil
}
