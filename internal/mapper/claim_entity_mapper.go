package mapper

import (
	"github.com/VulpesFerrilata/auth/internal/domain/entity"
	"github.com/VulpesFerrilata/auth/internal/domain/model"
)

type ClaimEntityMapper interface {
	ToClaimModel() *model.Claim
}

func NewClaimEntityMapper(claimEntity *entity.Claim) ClaimEntityMapper {
	return &claimEntityMapper{
		claimEntity: claimEntity,
	}
}

type claimEntityMapper struct {
	claimEntity *entity.Claim
}

func (c claimEntityMapper) ToClaimModel() *model.Claim {
	return model.ToClaim(
		c.claimEntity.ID,
		c.claimEntity.UserID,
		c.claimEntity.Jti,
		int64(c.claimEntity.Version),
	)
}
