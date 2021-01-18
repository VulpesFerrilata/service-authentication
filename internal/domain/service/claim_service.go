package service

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/app_error/authentication_error"
	"github.com/VulpesFerrilata/auth/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ClaimService interface {
	GetClaimRepository() repository.ClaimRepository
	ValidateAuthenticate(ctx context.Context, userId uuid.UUID, jti uuid.UUID) error
}

func NewClaimService(claimRepository repository.ClaimRepository) ClaimService {
	return &claimService{
		claimRepository: claimRepository,
	}
}

type claimService struct {
	claimRepository repository.ClaimRepository
}

func (cs claimService) GetClaimRepository() repository.ClaimRepository {
	return cs.claimRepository
}

func (cs claimService) ValidateAuthenticate(ctx context.Context, userId uuid.UUID, jti uuid.UUID) error {
	claim, err := cs.claimRepository.GetByUserId(ctx, userId)
	if err != nil {
		return errors.Wrap(err, "service.ClaimService.ValidateAuthenticate")
	}

	if claim.GetJti() != jti {
		return authentication_error.NewLoggedInByAnotherDeviceError()
	}

	return nil
}
