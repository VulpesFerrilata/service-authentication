package model

import "github.com/google/uuid"

func NewUser(id uuid.UUID, displayName string, email string) *User {
	user := new(User)
	user.id = id
	user.displayName = displayName
	user.email = email
	return user
}

type User struct {
	id          uuid.UUID
	displayName string
	email       string
}

func (u User) GetID() uuid.UUID {
	return u.id
}

func (u User) GetDisplayName() string {
	return u.displayName
}

func (u User) GetEmail() string {
	return u.email
}
