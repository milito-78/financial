package domain

import "time"

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	FirstName string
	LastName  string
	Uuid      string
	ID        uint64
}

func NewUser(username string, firstName string, lastName string, uuid string) *User {
	return &User{Username: username, FirstName: firstName, LastName: lastName, Uuid: uuid}
}
