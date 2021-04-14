package model

import (
	"github.com/VulpesFerrilata/auth/internal/pkg/app_error/authentication_error"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func NewClaim(userId uuid.UUID) (*Claim, error) {
	claim := new(Claim)
	claim.userID = userId

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	claim.id = id

	claim.userID = userId

	jti, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	claim.jti = jti

	return claim, nil
}

func ToClaim(id uuid.UUID, userId uuid.UUID, jti uuid.UUID, version int64) *Claim {
	claim := new(Claim)
	claim.id = id
	claim.userID = userId
	claim.jti = jti
	claim.version = version
	return claim
}

type Claim struct {
	id      uuid.UUID
	userID  uuid.UUID
	jti     uuid.UUID
	version int64
}

func (c Claim) GetId() uuid.UUID {
	return c.id
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

func (c Claim) ValidateJti(jti uuid.UUID) error {
	if c.jti != jti {
		return authentication_error.NewLoggedInByAnotherDeviceError()
	}
	return nil
}
