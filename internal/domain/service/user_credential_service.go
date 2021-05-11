package service

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/mapper"
	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/auth/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UserCredentialService interface {
	NewUserCredential(ctx context.Context, id uuid.UUID, username string, hashPassword []byte) (*model.UserCredential, error)
	GetByUsername(ctx context.Context, username string) (*model.UserCredential, error)
	Save(ctx context.Context, userCredential *model.UserCredential) (*model.UserCredential, error)
}

func NewUserCredentialService(userCredentialRepository repository.UserCredentialRepository,
	userCredentialMapper mapper.UserCredentialMapper) UserCredentialService {
	return &userCredentialService{
		userCredentialRepository: userCredentialRepository,
		userCredentialMapper:     userCredentialMapper,
	}
}

type userCredentialService struct {
	userCredentialRepository repository.UserCredentialRepository
	userCredentialMapper     mapper.UserCredentialMapper
}

func (c userCredentialService) NewUserCredential(ctx context.Context, id uuid.UUID, username string, hashPassword []byte) (*model.UserCredential, error) {
	userCredential := model.NewUserCredential(id, username, hashPassword)
	return userCredential, nil
}

func (c userCredentialService) GetByUsername(ctx context.Context, username string) (*model.UserCredential, error) {
	userCredentialEntity, err := c.userCredentialRepository.GetByUsername(ctx, username)

	return c.userCredentialMapper.GetModel(ctx, userCredentialEntity), errors.WithStack(err)
}

func (c userCredentialService) Save(ctx context.Context, userCredential *model.UserCredential) (*model.UserCredential, error) {
	userCredentialEntity := c.userCredentialMapper.GetEntity(ctx, userCredential)

	modelState := c.userCredentialMapper.GetModelState(ctx, userCredential)
	switch modelState {
	case mapper.New:
		if err := c.userCredentialRepository.Insert(ctx, userCredentialEntity); err != nil {
			return nil, errors.WithStack(err)
		}
	case mapper.Modified:
		if err := c.userCredentialRepository.Update(ctx, userCredentialEntity); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return c.userCredentialMapper.GetModel(ctx, userCredentialEntity), nil
}
