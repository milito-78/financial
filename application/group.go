package application

import (
	"errors"
	"financial/domain"
	"financial/infrastructure/db"
)

type IGroupService interface {
	UserGroupPaginate(user uint64, page uint) *db.Paginate[domain.Group]
	Store(group *domain.Group) error
	Update(group *domain.Group) error
	GetGroup(id uint64) (*domain.Group, error)
	DeleteGroup(id uint64) bool
	LeaveGroup(groupId uint64, userId uint64) bool
}

type GroupService struct {
	reader db.GroupReader
	writer db.GroupWriter
}

func NewGroupService(reader db.GroupReader, writer db.GroupWriter) *GroupService {
	return &GroupService{reader: reader, writer: writer}
}

func (g GroupService) UserGroupPaginate(user uint64, page uint) *db.Paginate[domain.Group] {
	return g.reader.UserGroupsPaginate(user, page)
}

func (g GroupService) Store(group *domain.Group) error {
	return g.writer.Create(group)
}

func (g GroupService) Update(group *domain.Group) error {
	return g.writer.Update(group)
}

func (g GroupService) GetGroup(id uint64) (*domain.Group, error) {
	group := g.reader.Get(id)
	if group == nil {
		return nil, errors.New("group not exists")
	}
	return group, nil
}

func (g GroupService) DeleteGroup(id uint64) bool {
	err := g.writer.SoftDelete(id)
	if err != nil {
		return false
	}
	return true
}

func (g GroupService) LeaveGroup(groupId uint64, userId uint64) bool {
	err := g.writer.UserLeaveGroup(groupId, userId)
	if err != nil {
		return false
	}
	return true
}
