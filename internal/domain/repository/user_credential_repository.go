package repository

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/entity"
	"github.com/VulpesFerrilata/library/pkg/app_error"
	"github.com/VulpesFerrilata/library/pkg/middleware"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SafeUserCredentialRepository interface {
	GetByUsername(ctx context.Context, username string) (*entity.UserCredential, error)
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

func (c userCredentialRepository) GetByUsername(ctx context.Context, username string) (*entity.UserCredential, error) {
	userCredentialEntity := new(entity.UserCredential)

	err := c.transactionMiddleware.Get(ctx).Where("username = ?", username).First(userCredentialEntity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, app_error.NewRecordNotFoundError("user credential")
	}
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
