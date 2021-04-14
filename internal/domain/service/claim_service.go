package service

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/auth/internal/domain/repository"
	"github.com/VulpesFerrilata/auth/internal/mapper"
	"github.com/VulpesFerrilata/library/pkg/app_error"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ClaimService interface {
	GetByUserId(ctx context.Context, userId uuid.UUID) (*model.Claim, error)
	Create(ctx context.Context, claim *model.Claim) (*model.Claim, error)
}

func NewClaimService(claimRepository repository.ClaimRepository) ClaimService {
	return &claimService{
		claimRepository: claimRepository,
	}
}

type claimService struct {
	claimRepository repository.ClaimRepository
}

func (c claimService) GetByUserId(ctx context.Context, userId uuid.UUID) (*model.Claim, error) {
	claimEntity, err := c.claimRepository.GetByUserId(ctx, userId)

	return mapper.NewClaimEntityMapper(claimEntity).ToClaimModel(), errors.WithStack(err)
}

func (c claimService) validate(ctx context.Context, claim *model.Claim) error {
	businessRuleErrs := make([]app_error.BusinessRuleError, 0)

	isExist, err := c.claimRepository.IsUserIdExists(ctx, claim.GetId(), claim.GetUserId())
	if err != nil {
		return errors.WithStack(err)
	}
	if isExist {
		appErr := app_error.NewAlreadyExistsError("user id")
		businessRuleErrs = append(businessRuleErrs, appErr)
	}

	if len(businessRuleErrs) > 0 {
		return app_error.NewBusinessRuleErrors(businessRuleErrs...)
	}

	return nil
}

func (c claimService) Create(ctx context.Context, claim *model.Claim) (*model.Claim, error) {
	if err := c.validate(ctx, claim); err != nil {
		return nil, errors.WithStack(err)
	}

	claimEntity := mapper.NewClaimModelMapper(claim).ToClaimEntity()
	if err := c.claimRepository.Insert(ctx, claimEntity); err != nil {
		return nil, errors.WithStack(err)
	}

	return mapper.NewClaimEntityMapper(claimEntity).ToClaimModel(), nil
}
