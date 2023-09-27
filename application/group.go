package application

import (
	"financial/domain"
	"financial/infrastructure/db"
)

type IGroupService interface {
	UserGroupPaginate(user uint64, page uint) db.Paginate[domain.Group]
	Create(group *domain.Group) error
	Update(group *domain.Group) error
	Get(id uint64) (*domain.Group, error)
}

type GroupService struct {
	reader db.GroupReader
	writer db.GroupWriter
}

func NewGroupService(reader db.GroupReader, writer db.GroupWriter) *GroupService {
	return &GroupService{reader: reader, writer: writer}
}

func (g GroupService) UserGroupPaginate(user uint64, page uint) db.Paginate[domain.Group] {
	return g.reader.UserGroupsPaginate(user, page)
}

func (g GroupService) Create(group *domain.Group) error {
	//TODO implement me
	panic("implement me")
}

func (g GroupService) Update(group *domain.Group) error {
	//TODO implement me
	panic("implement me")
}

func (g GroupService) Get(id uint64) (*domain.Group, error) {
	//TODO implement me
	panic("implement me")
}
