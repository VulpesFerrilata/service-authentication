package model

import (
	"github.com/google/uuid"
)

func NewClaim(id uuid.UUID, jti uuid.UUID) *Claim {
	claim := new(Claim)
	claim.id = id
	claim.jti = jti
	return claim
}

type Claim struct {
	id  uuid.UUID
	jti uuid.UUID
}

func (c Claim) GetId() uuid.UUID {
	return c.id
}

func (c Claim) GetJti() uuid.UUID {
	return c.jti
}

func (c *Claim) SetJti(jti uuid.UUID) {
	c.jti = jti
}
