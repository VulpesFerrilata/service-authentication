package model

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func NewClaim(userId uuid.UUID) (*Claim, error) {
	claim := new(Claim)
	claim.userID = userId

	jti, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	claim.jti = jti

	return claim, nil
}

func ToClaim(userId uuid.UUID, jti uuid.UUID, version int64) *Claim {
	claim := new(Claim)
	claim.userID = userId
	claim.jti = jti
	claim.version = version
	return claim
}

type Claim struct {
	userID  uuid.UUID
	jti     uuid.UUID
	version int64
}

func (c Claim) GetUserId() uuid.UUID {
	return c.userID
}

func (c Claim) GetJti() uuid.UUID {
	return c.jti
}

func (c Claim) GetVersion() int64 {
	return c.version
}
