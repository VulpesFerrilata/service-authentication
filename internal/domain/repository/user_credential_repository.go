package repository

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/entity"
	"github.com/VulpesFerrilata/library/pkg/middleware"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type SafeUserCredentialRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserCredential, error)
}

type UserCredentialRepository interface {
	SafeUserCredentialRepository
	Insert(ctx context.Context, userCredentialEntity *entity.UserCredential) error
	Update(ctx context.Context, userCredentialEntity *entity.UserCredential) error
}

func NewUserCredentialRepository(transactionMiddleware *middleware.TransactionMiddleware) UserCredentialRepository {
	return &userCredentialRepository{
		transactionMiddleware: transactionMiddleware,
	}
}

type userCredentialRepository struct {
	transactionMiddleware *middleware.TransactionMiddleware
}

func (c userCredentialRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserCredential, error) {
	userCredentialEntity := new(entity.UserCredential)
	err := c.transactionMiddleware.Get(ctx).Where("user_id = ?", userID).First(userCredentialEntity).Error
	return userCredentialEntity, errors.WithStack(err)
}

func (c userCredentialRepository) Insert(ctx context.Context, userCredentialEntity *entity.UserCredential) error {
	err := c.transactionMiddleware.Get(ctx).Create(userCredentialEntity).Error
	return errors.WithStack(err)
}

func (c userCredentialRepository) Update(ctx context.Context, userCredentialEntity *entity.UserCredential) error {
	err := c.transactionMiddleware.Get(ctx).Updates(userCredentialEntity).Error
	return errors.WithStack(err)
}
