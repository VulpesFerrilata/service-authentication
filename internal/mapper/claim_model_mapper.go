package mapper

import (
	"github.com/VulpesFerrilata/auth/internal/domain/entity"
	"github.com/VulpesFerrilata/auth/internal/domain/model"
	gorm_custom "github.com/VulpesFerrilata/library/pkg/gorm"
	"github.com/dgrijalva/jwt-go"
)

type ClaimModelMapper interface {
	ToClaimEntity() *entity.Claim
	ToStandardClaim() *jwt.StandardClaims
}

func NewClaimModelMapper(claim *model.Claim) ClaimModelMapper {
	return &claimModelMapper{
		claim: claim,
	}
}

type claimModelMapper struct {
	claim *model.Claim
}

func (c claimModelMapper) ToClaimEntity() *entity.Claim {
	claimEntity := new(entity.Claim)
	claimEntity.ID = c.claim.GetId()
	claimEntity.UserID = c.claim.GetUserId()
	claimEntity.Jti = c.claim.GetJti()
	claimEntity.Version = gorm_custom.Version(c.claim.GetVersion())
	return claimEntity
}

func (c claimModelMapper) ToStandardClaim() *jwt.StandardClaims {
	standardClaim := new(jwt.StandardClaims)
	standardClaim.Id = c.claim.GetJti().String()
	standardClaim.Subject = c.claim.GetUserId().String()
	return standardClaim
}
