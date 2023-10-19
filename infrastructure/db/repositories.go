package db

import "financial/domain"

type Paginate[T interface{}] struct {
	Results  []*T
	Page     uint
	NextPage bool
}

func MakeSimplePaginate[T interface{}](result []*T, page int, perPage int) *Paginate[T] {
	nextPage := false
	if len(result) > perPage {
		nextPage = true
		result = result[0:perPage]
	}

	return &Paginate[T]{
		Results:  result,
		Page:     uint(page),
		NextPage: nextPage,
	}
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
	UserGroupsPaginate(user uint64, page uint) *Paginate[domain.Group]
}

type GroupWriter interface {
	Create(group *domain.Group) error
	Update(group *domain.Group) error
	SoftDelete(id uint64) error
	UserLeaveGroup(groupId uint64, userId uint64) error
}
