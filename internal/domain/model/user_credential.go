package model

import "github.com/google/uuid"

func NewUserCredential(id uuid.UUID, username string, hashPassword []byte) *UserCredential {
	userCredential := new(UserCredential)
	userCredential.id = id
	userCredential.username = username
	userCredential.hashPassword = hashPassword
	return userCredential
}

type UserCredential struct {
	id           uuid.UUID
	username     string
	hashPassword []byte
}

func (u UserCredential) GetId() uuid.UUID {
	return u.id
}

func (u UserCredential) GetUsername() string {
	return u.username
}

func (u UserCredential) GetHashPassword() []byte {
	return u.hashPassword
}
