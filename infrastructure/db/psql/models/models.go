package models

import (
	"financial/domain"
	"time"
)

type Identifier struct {
	ID uint64 `gorm:"primaryKey"`
}

type Dates struct {
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:nano"`
}

type UserEntity struct {
	FirstName string
	Username  string `gorm:"size:255;index:idx_username,unique;not null"`
	LastName  string
	Uuid      string `gorm:"size:255;index:idx_uuid,unique;not null"`
	Identifier
	Dates
}

func (u *UserEntity) ToUser() *domain.User {
	return &domain.User{
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Uuid:      u.Uuid,
		ID:        u.ID,
	}
}

func FromUser(user *domain.User) *UserEntity {
	return &UserEntity{
		FirstName: user.FirstName,
		Username:  user.Username,
		LastName:  user.LastName,
		Uuid:      user.Uuid,
		Dates: Dates{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}
