package datamodel

import (
	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func NewClaim(userId uuid.UUID) (*Claim, error) {
	claim := new(Claim)
	claim.userID = userId

	jti, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(err, "datamodel.NewClaim")
	}
	claim.jti = jti

	return claim, nil
}

func NewClaimFromStandardClaim(standardClaim *jwt.StandardClaims) (*Claim, error) {
	claim := new(Claim)

	userId, err := uuid.Parse(standardClaim.Subject)
	if err != nil {
		return nil, errors.Wrap(err, "datamodel.NewClaimFromStandardClaim")
	}
	claim.userID = userId

	jti, err := uuid.Parse(standardClaim.Id)
	if err != nil {
		return nil, errors.Wrap(err, "datamodel.NewClaimFromStandardClaim")
	}
	claim.jti = jti
	return claim, nil
}

func NewClaimFromClaimModel(claimModel *model.Claim) *Claim {
	claim := new(Claim)
	claim.userID = claimModel.UserID
	claim.jti = claimModel.Jti
	return claim
}

type Claim struct {
	userID uuid.UUID
	jti    uuid.UUID
}

func (c Claim) GetUserId() uuid.UUID {
	return c.userID
}

func (c Claim) GetJti() uuid.UUID {
	return c.jti
}

func (c Claim) ToModel() *model.Claim {
	claimModel := new(model.Claim)
	claimModel.UserID = c.userID
	claimModel.Jti = c.jti
	return claimModel
}
