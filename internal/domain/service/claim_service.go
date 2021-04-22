package service

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/auth/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ClaimService interface {
	GetByUserId(ctx context.Context, userId uuid.UUID) (*model.Claim, error)
	Save(ctx context.Context, claim *model.Claim) (*model.Claim, error)
}

func NewClaimService(claimRepository repository.ClaimRepository) ClaimService {
	return &claimService{
		claimRepository:           claimRepository,
		claimChangeTrackerService: NewClaimChangeTrackerService(),
	}
}

type claimService struct {
	claimRepository           repository.ClaimRepository
	claimChangeTrackerService ClaimChangeTrackerService
}

func (c claimService) GetByUserId(ctx context.Context, userId uuid.UUID) (*model.Claim, error) {
	claimEntity, err := c.claimRepository.GetByUserId(ctx, userId)

	return c.claimChangeTrackerService.GetModel(ctx, claimEntity), errors.WithStack(err)
}

func (c claimService) Save(ctx context.Context, claim *model.Claim) (*model.Claim, error) {
	claimEntity := c.claimChangeTrackerService.GetEntity(ctx, claim)

	entityState := c.claimChangeTrackerService.GetEntityState(ctx, claim)
	switch entityState {
	case New:
		if err := c.claimRepository.Insert(ctx, claimEntity); err != nil {
			return nil, errors.WithStack(err)
		}
	case Modified:
		if err := c.claimRepository.Update(ctx, claimEntity); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return c.claimChangeTrackerService.GetModel(ctx, claimEntity), nil
}
