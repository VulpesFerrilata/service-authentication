package mapper

import (
	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/grpc/protoc/user"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UserResponsePbMapper interface {
	ToUserModel() (*model.User, error)
}

func NewUserResponsePbMapper(userResponsePb *user.UserResponse) UserResponsePbMapper {
	return &userResponsePbMapper{
		userResponsePb: userResponsePb,
	}
}

type userResponsePbMapper struct {
	userResponsePb *user.UserResponse
}

func (u userResponsePbMapper) ToUserModel() (*model.User, error) {
	id, err := uuid.Parse(u.userResponsePb.GetID())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return model.ToUser(
		id,
		u.userResponsePb.GetUsername(),
		u.userResponsePb.GetDisplayName(),
		u.userResponsePb.GetEmail(),
	), nil
}
