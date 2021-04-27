package service

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/mapper"
	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/auth/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ClaimService interface {
	GetById(ctx context.Context, id uuid.UUID) (*model.Claim, error)
	Save(ctx context.Context, claim *model.Claim) (*model.Claim, error)
}

func NewClaimService(claimRepository repository.ClaimRepository) ClaimService {
	return &claimService{
		claimRepository: claimRepository,
		claimMapper:     mapper.NewClaimMapper(),
	}
}

type claimService struct {
	claimRepository repository.ClaimRepository
	claimMapper     mapper.ClaimMapper
}

func (c claimService) GetById(ctx context.Context, id uuid.UUID) (*model.Claim, error) {
	claimEntity, err := c.claimRepository.GetById(ctx, id)

	return c.claimMapper.GetModel(ctx, claimEntity), errors.WithStack(err)
}

func (c claimService) Save(ctx context.Context, claim *model.Claim) (*model.Claim, error) {
	claimEntity := c.claimMapper.GetEntity(ctx, claim)

	modelState := c.claimMapper.GetModelState(ctx, claim)
	switch modelState {
	case mapper.New:
		if err := c.claimRepository.Insert(ctx, claimEntity); err != nil {
			return nil, errors.WithStack(err)
		}
	case mapper.Modified:
		if err := c.claimRepository.Update(ctx, claimEntity); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return c.claimMapper.GetModel(ctx, claimEntity), nil
}
