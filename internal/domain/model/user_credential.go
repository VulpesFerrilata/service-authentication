package model

import (
	"time"

	"github.com/google/uuid"
)

func NewUserCredential(id uuid.UUID, userID uuid.UUID, hashPassword []byte, createdAt time.Time) *UserCredential {
	userCredential := new(UserCredential)
	userCredential.id = id
	userCredential.userID = userID
	userCredential.hashPassword = hashPassword
	return userCredential
}

type UserCredential struct {
	id           uuid.UUID
	userID       uuid.UUID
	hashPassword []byte
	createdAt    time.Time
}

func (u UserCredential) GetID() uuid.UUID {
	return u.id
}

func (u UserCredential) GetUserID() uuid.UUID {
	return u.userID
}

func (u UserCredential) GetHashPassword() []byte {
	return u.hashPassword
}

func (u UserCredential) GetCreatedAt() time.Time {
	return u.createdAt
}
