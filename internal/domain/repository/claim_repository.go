package repository

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/entity"
	"github.com/VulpesFerrilata/library/pkg/app_error"
	"github.com/VulpesFerrilata/library/pkg/middleware"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SafeClaimRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (*entity.Claim, error)
}

type ClaimRepository interface {
	SafeClaimRepository
	Insert(ctx context.Context, claimEntity *entity.Claim) error
	Update(ctx context.Context, claimEntity *entity.Claim) error
}

func NewClaimRepository(transactionMiddleware *middleware.TransactionMiddleware) ClaimRepository {
	return &claimRepository{
		transactionMiddleware: transactionMiddleware,
	}
}

type claimRepository struct {
	transactionMiddleware *middleware.TransactionMiddleware
}

func (c claimRepository) GetById(ctx context.Context, id uuid.UUID) (*entity.Claim, error) {
	claimEntity := new(entity.Claim)

	err := c.transactionMiddleware.Get(ctx).First(claimEntity, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = app_error.NewRecordNotFoundError("claim")
	}
	return claimEntity, errors.WithStack(err)
}

func (c claimRepository) Insert(ctx context.Context, claimEntity *entity.Claim) error {
	err := c.transactionMiddleware.Get(ctx).Create(claimEntity).Error
	return errors.WithStack(err)
}

func (c claimRepository) Update(ctx context.Context, claimEntity *entity.Claim) error {
	err := c.transactionMiddleware.Get(ctx).Updates(claimEntity).Error
	return errors.WithStack(err)
}
