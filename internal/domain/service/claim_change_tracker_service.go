package service

import (
	"context"
	"reflect"
	"sync"

	"github.com/VulpesFerrilata/auth/internal/domain/entity"
	"github.com/VulpesFerrilata/auth/internal/domain/model"
)

type entityState int

const (
	New entityState = iota
	Modified
	Unchanged
)

type ClaimChangeTrackerService interface {
	GetEntityState(ctx context.Context, claim *model.Claim) entityState
	GetModel(ctx context.Context, claimEntity *entity.Claim) *model.Claim
	GetEntity(ctx context.Context, claim *model.Claim) *entity.Claim
}

func NewClaimChangeTrackerService() ClaimChangeTrackerService {
	return &claimChangeTrackerService{
		claimMap: make(map[*model.Claim]entity.Claim),
	}
}

type claimChangeTrackerService struct {
	claimMap     map[*model.Claim]entity.Claim
	claimMapLock sync.RWMutex
}

func (u *claimChangeTrackerService) GetModel(ctx context.Context, claimEntity *entity.Claim) *model.Claim {
	claim := model.NewClaim(
		claimEntity.UserID,
		claimEntity.Jti,
	)

	u.claimMapLock.Lock()
	u.claimMap[claim] = *claimEntity
	u.claimMapLock.Unlock()

	go func(done <-chan struct{}) {
		<-done
		u.claimMapLock.Lock()
		delete(u.claimMap, claim)
		u.claimMapLock.Unlock()
	}(ctx.Done())

	return claim
}

func (u *claimChangeTrackerService) GetEntity(ctx context.Context, claim *model.Claim) *entity.Claim {
	newClaimEntity := new(entity.Claim)
	newClaimEntity.UserID = claim.GetUserId()
	newClaimEntity.Jti = claim.GetJti()

	u.claimMapLock.RLock()
	claimEntity, ok := u.claimMap[claim]
	u.claimMapLock.RUnlock()
	if ok {
		newClaimEntity.Version = claimEntity.Version
	}

	return newClaimEntity
}

func (u *claimChangeTrackerService) GetEntityState(ctx context.Context, claim *model.Claim) entityState {
	u.claimMapLock.RLock()
	claimEntity, ok := u.claimMap[claim]
	u.claimMapLock.RUnlock()
	if !ok {
		return New
	}

	newClaimEntity := u.GetEntity(ctx, claim)

	if reflect.DeepEqual(&claimEntity, newClaimEntity) {
		return Unchanged
	}

	return Modified
}
