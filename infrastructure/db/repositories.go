package db

import "financial/domain"

type NotFoundError struct {
}

func (n NotFoundError) Error() string {
	return "not found"
}

type UserReader interface {
	Get(id uint64) *domain.User
	GetByUuid(uid string) *domain.User
	GetByUsername(username string) *domain.User
}

type UserWriter interface {
	Create(user *domain.User) error
	Update(user *domain.User) error
}
