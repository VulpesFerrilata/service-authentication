package repository

import (
	"context"
	"errors"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/VulpesFerrilata/library/pkg/db"
	server_errors "github.com/VulpesFerrilata/library/pkg/errors"
	"gorm.io/gorm"
)

type SafeClaimRepository interface {
	GetByUserId(ctx context.Context, userId uint) (*model.Claim, error)
}

type ClaimRepository interface {
	SafeClaimRepository
	Insert(ctx context.Context, claim *model.Claim) error
	DeleteByUserId(ctx context.Context, claim *model.Claim) error
}

func NewClaimRepository(dbContext *db.DbContext) ClaimRepository {
	return &claimRepository{
		dbContext: dbContext,
	}
}

type claimRepository struct {
	dbContext *db.DbContext
}

func (tr claimRepository) GetByUserId(ctx context.Context, userId uint) (*model.Claim, error) {
	claim := new(model.Claim)
	err := tr.dbContext.GetDB(ctx).First(claim, "user_id = ?", userId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = server_errors.NewNotFoundError("claim")
	}
	return claim, err
}

func (tr claimRepository) Insert(ctx context.Context, claim *model.Claim) error {
	return tr.dbContext.GetDB(ctx).Create(claim).Error
}

func (tr claimRepository) DeleteByUserId(ctx context.Context, claim *model.Claim) error {
	return tr.dbContext.GetDB(ctx).Delete(claim, "user_id = ?", claim.UserID).Error
}
