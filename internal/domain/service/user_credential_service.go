package service

import (
	"context"
	"time"

	"github.com/VulpesFerrilata/auth/internal/domain/mapper"
	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/auth/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UserCredentialService interface {
	NewUserCredential(ctx context.Context, userID uuid.UUID, hashPassword []byte) (*model.UserCredential, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.UserCredential, error)
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

func (u userCredentialService) NewUserCredential(ctx context.Context, userID uuid.UUID, hashPassword []byte) (*model.UserCredential, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userCredential := model.NewUserCredential(id, userID, hashPassword, time.Now())
	return userCredential, nil
}

func (u userCredentialService) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.UserCredential, error) {
	userCredentialEntity, err := u.userCredentialRepository.GetByUserID(ctx, userID)
	return u.userCredentialMapper.GetModel(ctx, userCredentialEntity), errors.WithStack(err)
}

func (u userCredentialService) Save(ctx context.Context, userCredential *model.UserCredential) (*model.UserCredential, error) {
	userCredentialEntity := u.userCredentialMapper.GetEntity(ctx, userCredential)

	modelState := u.userCredentialMapper.GetModelState(ctx, userCredential)
	switch modelState {
	case mapper.New:
		if err := u.userCredentialRepository.Insert(ctx, userCredentialEntity); err != nil {
			return nil, errors.WithStack(err)
		}
	case mapper.Modified:
		if err := u.userCredentialRepository.Update(ctx, userCredentialEntity); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return u.userCredentialMapper.GetModel(ctx, userCredentialEntity), nil
}
