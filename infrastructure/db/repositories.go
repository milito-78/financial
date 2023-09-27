package db

import "financial/domain"

type Paginate[T interface{}] struct {
	Results  []*T
	Page     uint
	NextPage bool
}

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

type GroupReader interface {
	Get(id uint64) *domain.Group
	UserGroupsPaginate(user uint64, page uint) Paginate[domain.Group]
}

type GroupWriter interface {
	Create(group *domain.Group) error
	Update(group *domain.Group) error
}
