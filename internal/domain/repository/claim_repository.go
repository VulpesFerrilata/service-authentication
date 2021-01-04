package repository

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/datamodel"
	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/library/pkg/app_error"
	"github.com/VulpesFerrilata/library/pkg/middleware"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
)

type SafeClaimRepository interface {
	GetByUserId(ctx context.Context, userId int) (*datamodel.Claim, error)
}

type ClaimRepository interface {
	SafeClaimRepository
	InsertOrUpdate(ctx context.Context, claim *datamodel.Claim) error
}

func NewClaimRepository(transactionMiddleware *middleware.TransactionMiddleware,
	validate *validator.Validate) ClaimRepository {
	return &claimRepository{
		transactionMiddleware: transactionMiddleware,
		validate:              validate,
	}
}

type claimRepository struct {
	transactionMiddleware *middleware.TransactionMiddleware
	validate              *validator.Validate
}

func (tr claimRepository) GetByUserId(ctx context.Context, userId int) (*datamodel.Claim, error) {
	claimModel := new(model.Claim)

	err := tr.transactionMiddleware.Get(ctx).First(claimModel, userId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = app_error.NewNotFoundError("claim")
	}
	return datamodel.NewClaimFromClaimModel(claimModel), errors.Wrap(err, "repository.ClaimRepository.GetByUserId")
}

func (tr claimRepository) InsertOrUpdate(ctx context.Context, claim *datamodel.Claim) error {
	return claim.Persist(func(claimModel *model.Claim) error {
		if err := tr.validate.StructCtx(ctx, claimModel); err != nil {
			if fieldErrors, ok := errors.Cause(err).(validator.ValidationErrors); ok {
				err = app_error.NewEntityValidationError(claimModel, fieldErrors)
			}
			return errors.Wrap(err, "repository.ClaimRepository.InsertOrUpdate")
		}

		err := tr.transactionMiddleware.Get(ctx).Save(claimModel).Error
		return errors.Wrap(err, "repository.ClaimRepository.InsertOrUpdate")
	})
}
