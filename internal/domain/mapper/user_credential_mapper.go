package mapper

import (
	"context"
	"reflect"
	"sync"

	"github.com/VulpesFerrilata/auth/internal/domain/entity"
	"github.com/VulpesFerrilata/auth/internal/domain/model"
)

type UserCredentialMapper interface {
	GetModelState(ctx context.Context, userCredential *model.UserCredential) modelState
	GetModel(ctx context.Context, userCredentialEntity *entity.UserCredential) *model.UserCredential
	GetEntity(ctx context.Context, userCredential *model.UserCredential) *entity.UserCredential
}

func NewUserCredentialMapper() UserCredentialMapper {
	return &userCredentialMapper{
		userCredentialMap: make(map[*model.UserCredential]entity.UserCredential),
	}
}

type userCredentialMapper struct {
	userCredentialMap     map[*model.UserCredential]entity.UserCredential
	userCredentialMapLock sync.RWMutex
}

func (c *userCredentialMapper) GetModelState(ctx context.Context, userCredential *model.UserCredential) modelState {
	c.userCredentialMapLock.RLock()
	userCredentialEntity, ok := c.userCredentialMap[userCredential]
	c.userCredentialMapLock.RUnlock()
	if !ok {
		return New
	}

	newUserCredentialEntity := c.GetEntity(ctx, userCredential)

	if !reflect.DeepEqual(&userCredentialEntity, newUserCredentialEntity) {
		return Modified
	}

	return Unchanged
}

func (c *userCredentialMapper) GetModel(ctx context.Context, userCredentialEntity *entity.UserCredential) *model.UserCredential {
	if userCredentialEntity == nil {
		return nil
	}

	userCredential := model.NewUserCredential(
		userCredentialEntity.ID,
		userCredentialEntity.Username,
		userCredentialEntity.HashPassword,
	)

	c.userCredentialMapLock.Lock()
	c.userCredentialMap[userCredential] = *userCredentialEntity
	c.userCredentialMapLock.Unlock()

	go func(done <-chan struct{}) {
		<-done
		c.userCredentialMapLock.Lock()
		delete(c.userCredentialMap, userCredential)
		c.userCredentialMapLock.Unlock()
	}(ctx.Done())

	return userCredential
}

func (c *userCredentialMapper) GetEntity(ctx context.Context, userCredential *model.UserCredential) *entity.UserCredential {
	if userCredential == nil {
		return nil
	}

	newUserCredentialEntity := new(entity.UserCredential)
	newUserCredentialEntity.ID = userCredential.GetId()
	newUserCredentialEntity.Username = userCredential.GetUsername()
	newUserCredentialEntity.HashPassword = userCredential.GetHashPassword()

	c.userCredentialMapLock.RLock()
	userCredentialEntity, ok := c.userCredentialMap[userCredential]
	c.userCredentialMapLock.RUnlock()
	if ok {
		newUserCredentialEntity.Version = userCredentialEntity.Version
	}

	return newUserCredentialEntity
}
