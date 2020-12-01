package service

import (
	"context"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/auth/internal/domain/repository"
	server_errors "github.com/VulpesFerrilata/library/pkg/errors"
	"github.com/VulpesFerrilata/library/pkg/middleware"
	"github.com/VulpesFerrilata/library/pkg/validator"
)

type ClaimService interface {
	GetClaimRepository() repository.SafeClaimRepository
	ValidateAuthenticate(ctx context.Context, claim *model.Claim) error
	Create(ctx context.Context, claim *model.Claim) error
}

func NewClaimService(validate validator.Validate,
	claimRepository repository.ClaimRepository,
	translatorMiddleware *middleware.TranslatorMiddleware) ClaimService {
	return &claimService{
		validate:             validate,
		claimRepository:      claimRepository,
		translatorMiddleware: translatorMiddleware,
	}
}

type claimService struct {
	validate             validator.Validate
	claimRepository      repository.ClaimRepository
	translatorMiddleware *middleware.TranslatorMiddleware
}

func (cs claimService) GetClaimRepository() repository.SafeClaimRepository {
	return cs.claimRepository
}

func (cs claimService) ValidateAuthenticate(ctx context.Context, claim *model.Claim) error {
	trans := cs.translatorMiddleware.Get(ctx)

	if err := cs.validate.Struct(ctx, claim.Claim); err != nil {
		return err
	}

	validationErrs := server_errors.NewValidationError()
	claimDB, err := cs.claimRepository.GetByUserId(ctx, claim.UserID)
	if err != nil {
		return err
	}

	if claim.Jti != claimDB.Jti {
		fieldErr, _ := trans.T("validation-duplicate-login")
		validationErrs.WithFieldError(fieldErr)
	}

	if validationErrs.HasErrors() {
		return validationErrs
	}
	return nil
}

func (cs claimService) Create(ctx context.Context, claim *model.Claim) error {
	if err := claim.Init(); err != nil {
		return err
	}

	if err := cs.validate.Struct(ctx, claim); err != nil {
		return err
	}

	if err := cs.claimRepository.DeleteByUserId(ctx, claim); err != nil {
		return err
	}

	return cs.claimRepository.Insert(ctx, claim)
}
