package datamodel

import (
	"strconv"

	"github.com/VulpesFerrilata/auth/internal/domain/model"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/iris-contrib/go.uuid"
	"github.com/pkg/errors"
)

func NewClaim(userId int) (*Claim, error) {
	claim := new(Claim)
	claim.userID = userId

	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "datamodel.NewClaim")
	}
	claim.jti = uuid.String()

	return claim, nil
}

func NewClaimFromStandardClaim(standardClaim *jwt.StandardClaims) (*Claim, error) {
	claim := new(Claim)
	userId, err := strconv.ParseInt(standardClaim.Subject, 10, 32)
	if err != nil {
		return nil, errors.Wrap(err, "datamodel.NewClaimFromStandardClaim")
	}
	claim.userID = int(userId)
	claim.jti = standardClaim.Id
	return claim, nil
}

func NewClaimFromClaimModel(claimModel *model.Claim) *Claim {
	claim := new(Claim)
	claim.userID = claimModel.UserID
	claim.jti = claimModel.Jti
	return claim
}

type Claim struct {
	userID int
	jti    string
}

func (c Claim) GetUserId() int {
	return c.userID
}

func (c Claim) GetJti() string {
	return c.jti
}

func (c *Claim) Persist(f func(claimModel *model.Claim) error) error {
	claimModel := new(model.Claim)
	claimModel.UserID = c.userID
	claimModel.Jti = c.jti

	if err := f(claimModel); err != nil {
		return errors.Wrap(err, "datamodel.Claim.Persist")
	}

	c.userID = claimModel.UserID
	c.jti = claimModel.Jti

	return nil
}
