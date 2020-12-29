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
	GetByUserId(ctx context.Context, userId uint) (*model.Claim, error)
}

type ClaimRepository interface {
	SafeClaimRepository
	InsertOrUpdate(ctx context.Context, claim *model.Claim) error
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

func (tr claimRepository) GetByUserId(ctx context.Context, userId uint) (*model.Claim, error) {
	claim := model.EmptyClaim()

	return claim, claim.Persist(func(claim *datamodel.Claim) error {
		err := tr.transactionMiddleware.Get(ctx).First(claim, userId).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = app_error.NewNotFoundError("claim")
		}
		return errors.Wrap(err, "repository.ClaimRepository.GetByUserId")
	})
}

func (tr claimRepository) InsertOrUpdate(ctx context.Context, claim *model.Claim) error {
	return claim.Persist(func(claim *datamodel.Claim) error {
		if err := tr.validate.StructCtx(ctx, claim); err != nil {
			if fieldErrors, ok := errors.Cause(err).(validator.ValidationErrors); ok {
				err = app_error.NewValidationError(app_error.EntityValidation, "claim", fieldErrors)
			}
			return errors.Wrap(err, "repository.ClaimRepository.InsertOrUpdate")
		}

		err := tr.transactionMiddleware.Get(ctx).Save(claim).Error
		return errors.Wrap(err, "repository.ClaimRepository.InsertOrUpdate")
	})
}
