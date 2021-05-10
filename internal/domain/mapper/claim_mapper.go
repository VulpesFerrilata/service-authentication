package mapper

import (
	"context"
	"reflect"
	"sync"

	"github.com/VulpesFerrilata/auth/internal/domain/entity"
	"github.com/VulpesFerrilata/auth/internal/domain/model"
)

type ClaimMapper interface {
	GetModelState(ctx context.Context, claim *model.Claim) modelState
	GetModel(ctx context.Context, claimEntity *entity.Claim) *model.Claim
	GetEntity(ctx context.Context, claim *model.Claim) *entity.Claim
}

func NewClaimMapper() ClaimMapper {
	return &claimMapper{
		claimMap: make(map[*model.Claim]entity.Claim),
	}
}

type claimMapper struct {
	claimMap     map[*model.Claim]entity.Claim
	claimMapLock sync.RWMutex
}

func (c *claimMapper) GetModelState(ctx context.Context, claim *model.Claim) modelState {
	c.claimMapLock.RLock()
	claimEntity, ok := c.claimMap[claim]
	c.claimMapLock.RUnlock()
	if !ok {
		return New
	}

	newClaimEntity := c.GetEntity(ctx, claim)

	if !reflect.DeepEqual(&claimEntity, newClaimEntity) {
		return Modified
	}

	return Unchanged
}

func (c *claimMapper) GetModel(ctx context.Context, claimEntity *entity.Claim) *model.Claim {
	if claimEntity == nil {
		return nil
	}

	claim := model.NewClaim(
		claimEntity.UserID,
		claimEntity.Jti,
	)

	c.claimMapLock.Lock()
	c.claimMap[claim] = *claimEntity
	c.claimMapLock.Unlock()

	go func(done <-chan struct{}) {
		<-done
		c.claimMapLock.Lock()
		delete(c.claimMap, claim)
		c.claimMapLock.Unlock()
	}(ctx.Done())

	return claim
}

func (c *claimMapper) GetEntity(ctx context.Context, claim *model.Claim) *entity.Claim {
	newClaimEntity := new(entity.Claim)
	newClaimEntity.UserID = claim.GetUserID()
	newClaimEntity.Jti = claim.GetJti()

	c.claimMapLock.RLock()
	claimEntity, ok := c.claimMap[claim]
	c.claimMapLock.RUnlock()
	if ok {
		newClaimEntity.Version = claimEntity.Version
	}

	return newClaimEntity
}
