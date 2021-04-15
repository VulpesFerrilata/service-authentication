package service

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/auth/internal/domain/repository"
	"github.com/VulpesFerrilata/auth/internal/mapper"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ClaimService interface {
	GetByUserId(ctx context.Context, userId uuid.UUID) (*model.Claim, error)
	Save(ctx context.Context, claim *model.Claim) (*model.Claim, error)
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

func (c claimService) Save(ctx context.Context, claim *model.Claim) (*model.Claim, error) {
	claimEntity := mapper.NewClaimModelMapper(claim).ToClaimEntity()
	isExists, err := c.claimRepository.IsUserIdExists(ctx, claimEntity.UserID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if isExists {
		if err := c.claimRepository.Update(ctx, claimEntity); err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		if err := c.claimRepository.Insert(ctx, claimEntity); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return mapper.NewClaimEntityMapper(claimEntity).ToClaimModel(), nil
}
