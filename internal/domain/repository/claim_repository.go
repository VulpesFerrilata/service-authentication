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
	IsUserIdExists(ctx context.Context, userId uuid.UUID) (bool, error)
	GetByUserId(ctx context.Context, userId uuid.UUID) (*entity.Claim, error)
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

func (c claimRepository) IsUserIdExists(ctx context.Context, userId uuid.UUID) (bool, error) {
	var count int64

	if err := c.transactionMiddleware.Get(ctx).Model(&entity.Claim{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return false, errors.WithStack(err)
	}

	if count != 0 {
		return true, nil
	}

	return false, nil
}

func (c claimRepository) GetByUserId(ctx context.Context, userId uuid.UUID) (*entity.Claim, error) {
	claimEntity := new(entity.Claim)

	err := c.transactionMiddleware.Get(ctx).Where("user_id = ?", userId).First(claimEntity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = app_error.NewNotFoundError("claim")
	}
	return claimEntity, errors.WithStack(err)
}

func (c claimRepository) Insert(ctx context.Context, claimEntity *entity.Claim) error {
	err := c.transactionMiddleware.Get(ctx).Create(claimEntity).Error
	return errors.WithStack(err)
}

func (c claimRepository) Update(ctx context.Context, claimEntity *entity.Claim) error {
	tx := c.transactionMiddleware.Get(ctx).Updates(claimEntity)
	if tx.RowsAffected == 0 {
		return UpdateStaleObjectErr
	}
	return errors.WithStack(tx.Error)
}
