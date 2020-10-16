package repository

import (
	"context"
	"errors"

	"github.com/VulpesFerrilata/auth/internal/domain/datamodel"
	"github.com/VulpesFerrilata/library/pkg/db"
	server_errors "github.com/VulpesFerrilata/library/pkg/errors"
	"gorm.io/gorm"
)

type SafeClaimRepository interface {
	GetByUserId(ctx context.Context, userId uint) (*datamodel.Claim, error)
	GetByJti(ctx context.Context, jti string) (*datamodel.Claim, error)
}

type ClaimRepository interface {
	SafeClaimRepository
	Insert(ctx context.Context, claim *datamodel.Claim) error
	SaveByUserId(ctx context.Context, claim *datamodel.Claim) (int, error)
}

func NewClaimRepository(dbContext *db.DbContext) ClaimRepository {
	return &claimRepository{
		dbContext: dbContext,
	}
}

type claimRepository struct {
	dbContext *db.DbContext
}

func (tr claimRepository) GetByUserId(ctx context.Context, userId uint) (*datamodel.Claim, error) {
	claim := new(datamodel.Claim)
	err := tr.dbContext.GetDB(ctx).First(claim, "user_id = ?", userId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = server_errors.NewNotFoundError("claim")
	}
	return claim, err
}

func (tr claimRepository) GetByJti(ctx context.Context, jti string) (*datamodel.Claim, error) {
	claim := new(datamodel.Claim)
	err := tr.dbContext.GetDB(ctx).First(claim, "jti = ?", jti).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = server_errors.NewNotFoundError("claim")
	}
	return claim, err
}

func (tr claimRepository) Insert(ctx context.Context, claim *datamodel.Claim) error {
	return tr.dbContext.GetDB(ctx).Model(claim).Create(claim).Error
}

func (tr claimRepository) SaveByUserId(ctx context.Context, claim *datamodel.Claim) (int, error) {
	tx := tr.dbContext.GetDB(ctx).Model(claim).Where("user_id = ?", claim.UserID).Updates(claim)
	return int(tx.RowsAffected), tx.Error
}
