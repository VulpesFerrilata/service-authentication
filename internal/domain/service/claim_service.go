package service

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/business_rule_error"
	"github.com/VulpesFerrilata/auth/internal/domain/datamodel"
	"github.com/VulpesFerrilata/auth/internal/domain/repository"
	"github.com/VulpesFerrilata/library/pkg/app_error"
	"github.com/pkg/errors"
)

type ClaimService interface {
	GetClaimRepository() repository.SafeClaimRepository
	ValidateAuthenticate(ctx context.Context, userId int, jti string) error
	Save(ctx context.Context, claim *datamodel.Claim) error
}

func NewClaimService(claimRepository repository.ClaimRepository) ClaimService {
	return &claimService{
		claimRepository: claimRepository,
	}
}

type claimService struct {
	claimRepository repository.ClaimRepository
}

func (cs claimService) GetClaimRepository() repository.SafeClaimRepository {
	return cs.claimRepository
}

func (cs claimService) ValidateAuthenticate(ctx context.Context, userId int, jti string) error {
	claim, err := cs.claimRepository.GetByUserId(ctx, userId)
	if err != nil {
		return errors.Wrap(err, "service.ClaimService.ValidateAuthenticate")
	}

	if claim.GetJti() != jti {
		var businessRuleErrors app_error.BusinessRuleErrors
		loggedInByAnotherDeviceError := business_rule_error.NewLoggedInByAnotherDeviceError()
		businessRuleErrors = append(businessRuleErrors, loggedInByAnotherDeviceError)
		return businessRuleErrors
	}

	return nil
}

func (cs claimService) Save(ctx context.Context, claim *datamodel.Claim) error {
	err := cs.claimRepository.InsertOrUpdate(ctx, claim)
	return errors.Wrap(err, "service.ClaimService.Save")
}
