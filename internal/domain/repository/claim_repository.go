package repository

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/datamodel"
	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/library/pkg/app_error"
	"github.com/VulpesFerrilata/library/pkg/middleware"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SafeClaimRepository interface {
	GetByUserId(ctx context.Context, userId uint) (*model.Claim, error)
}

type ClaimRepository interface {
	SafeClaimRepository
	Insert(ctx context.Context, claim *model.Claim) error
	DeleteByUserId(ctx context.Context, userId uint) error
}

func NewClaimRepository(transactionMiddleware *middleware.TransactionMiddleware) ClaimRepository {
	return &claimRepository{
		transactionMiddleware: transactionMiddleware,
	}
}

type claimRepository struct {
	transactionMiddleware *middleware.TransactionMiddleware
}

func (tr claimRepository) GetByUserId(ctx context.Context, userId uint) (*model.Claim, error) {
	claim := model.EmptyClaim()

	return claim, claim.Persist(func(claim *datamodel.Claim) error {
		err := tr.transactionMiddleware.Get(ctx).First(claim, "user_id = ?", userId).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = app_error.NewNotFoundError("claim")
		}
		return errors.Wrap(err, "repository.ClaimRepository.GetByUserId")
	})
}

func (tr claimRepository) Insert(ctx context.Context, claim *model.Claim) error {
	return claim.Persist(func(claim *datamodel.Claim) error {
		err := tr.transactionMiddleware.Get(ctx).Create(claim).Error
		return errors.Wrap(err, "repository.ClaimRepository.Insert")
	})
}

func (tr claimRepository) DeleteByUserId(ctx context.Context, userId uint) error {
	claim := new(datamodel.Claim)
	err := tr.transactionMiddleware.Get(ctx).Delete(claim, "user_id = ?", userId).Error
	return errors.Wrap(err, "repository.ClaimRepository.DeleteByUserId")
}
