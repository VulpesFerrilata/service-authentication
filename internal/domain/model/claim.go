package model

import (
	"github.com/VulpesFerrilata/auth/internal/domain/datamodel"
	uuid "github.com/iris-contrib/go.uuid"
	"github.com/pkg/errors"
)

func NewClaim(userId uint) (*Claim, error) {
	claim := new(Claim)
	claim.userID = userId

	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "model.NewClaim")
	}
	claim.jti = uuid.String()

	return claim, nil
}

func NewClaimWithJti(userId uint, jti string) *Claim {
	claim := new(Claim)
	claim.userID = userId
	claim.jti = jti
	return claim
}

func EmptyClaim() *Claim {
	return new(Claim)
}

type Claim struct {
	userID uint
	jti    string
}

func (c Claim) GetUserId() uint {
	return c.userID
}

func (c Claim) GetJti() string {
	return c.jti
}

func (c *Claim) Persist(f func(claim *datamodel.Claim) error) error {
	claim := new(datamodel.Claim)
	claim.UserID = c.userID
	claim.Jti = c.jti

	if err := f(claim); err != nil {
		return errors.Wrap(err, "model.Claim.Persist")
	}

	c.userID = claim.UserID
	c.jti = claim.Jti

	return nil
}
