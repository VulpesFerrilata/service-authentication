package service

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/datamodel"
	"github.com/VulpesFerrilata/auth/internal/domain/repository"
	server_errors "github.com/VulpesFerrilata/library/pkg/errors"
	"github.com/VulpesFerrilata/library/pkg/middleware"
	uuid "github.com/iris-contrib/go.uuid"
)

type ClaimService interface {
	GetClaimRepository() repository.SafeClaimRepository
	ValidateAuthenticate(ctx context.Context, claim *datamodel.Claim) error
	Create(ctx context.Context, claim *datamodel.Claim) error
}

func NewClaimService(claimRepository repository.ClaimRepository,
	translatorMiddleware *middleware.TranslatorMiddleware) ClaimService {
	return &claimService{
		claimRepository:      claimRepository,
		translatorMiddleware: translatorMiddleware,
	}
}

type claimService struct {
	claimRepository      repository.ClaimRepository
	translatorMiddleware *middleware.TranslatorMiddleware
}

func (cs claimService) GetClaimRepository() repository.SafeClaimRepository {
	return cs.claimRepository
}

func (cs claimService) ValidateAuthenticate(ctx context.Context, claim *datamodel.Claim) error {
	trans := cs.translatorMiddleware.Get(ctx)
	validationErrs := server_errors.NewValidationError()
	claimDB, err := cs.claimRepository.GetByUserId(ctx, claim.UserID)
	if err != nil {
		return err
	}

	if claim.Jti != claimDB.Jti {
		fieldErr, _ := trans.T("validation-invalid", "jti")
		validationErrs.WithFieldError(fieldErr)
	}

	if validationErrs.HasErrors() {
		return validationErrs
	}
	return nil
}

func (cs claimService) Create(ctx context.Context, claim *datamodel.Claim) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	claim.Jti = uuid.String()

	rowAffected, err := cs.claimRepository.SaveByUserId(ctx, claim)
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return cs.claimRepository.Insert(ctx, claim)
	}

	return nil
}
