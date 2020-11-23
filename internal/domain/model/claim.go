package model

import (
	"github.com/VulpesFerrilata/auth/internal/domain/datamodel"
	uuid "github.com/iris-contrib/go.uuid"
)

type Claim struct {
	datamodel.Claim
}

func (c *Claim) Init() error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	c.Jti = uuid.String()
	return nil
}
