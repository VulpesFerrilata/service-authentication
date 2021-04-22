package model

import (
	"github.com/google/uuid"
)

func NewClaim(userId uuid.UUID, jti uuid.UUID) *Claim {
	claim := new(Claim)
	claim.userID = userId
	claim.jti = jti
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

func (c *Claim) SetJti(jti uuid.UUID) {
	c.jti = jti
}
