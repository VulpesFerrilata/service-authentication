package service

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/auth/internal/domain/repository"
	"github.com/VulpesFerrilata/library/pkg/app_error"
	"github.com/pkg/errors"
)

type ClaimService interface {
	GetClaimRepository() repository.SafeClaimRepository
	GetOrNewClaim(ctx context.Context, userId uint) (*model.Claim, error)
	ValidateAuthenticate(ctx context.Context, userId uint, jti string) error
	Save(ctx context.Context, claim *model.Claim) error
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

func (cs claimService) GetOrNewClaim(ctx context.Context, userId uint) (*model.Claim, error) {
	claim, err := cs.claimRepository.GetByUserId(ctx, userId)
	if err != nil {
		if _, ok := errors.Cause(err).(*app_error.NotFoundError); ok {
			claim, err = model.NewClaim(userId)
			return claim, errors.Wrap(err, "service.ClaimService.NewClaim")
		}
		return nil, errors.Wrap(err, "service.ClaimService.NewClaim")
	}
	return claim, nil
}

func (cs claimService) ValidateAuthenticate(ctx context.Context, userId uint, jti string) error {
	claim, err := cs.claimRepository.GetByUserId(ctx, userId)
	if err != nil {
		return errors.Wrap(err, "service.ClaimService.ValidateAuthenticate")
	}

	if claim.GetJti() != jti {

	}

	return nil
}

func (cs claimService) Save(ctx context.Context, claim *model.Claim) error {
	if err := cs.claimRepository.DeleteByUserId(ctx, claim.GetUserId()); err != nil {
		return errors.Wrap(err, "service.ClaimService.Create")
	}
	err := cs.claimRepository.Insert(ctx, claim)
	return errors.Wrap(err, "service.ClaimService.Create")
}
