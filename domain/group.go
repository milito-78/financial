package domain

import "time"

type Group struct {
	DeletedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatorId uint64
	Creator   User
	Name      string
	ID        uint64
}

func NewGroup(creatorId uint64, name string, ID uint64) *Group {
	return &Group{CreatorId: creatorId, Name: name, ID: ID}
}
