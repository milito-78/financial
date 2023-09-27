package db

import (
	"errors"
	"financial/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u UserRepository) Create(user *domain.User) error {
	model := FromUser(user)
	res := u.db.Create(model)
	if res.Error != nil {
		return res.Error
	}
	user.ID = model.ID
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = model.UpdatedAt
	return nil
}

func (u UserRepository) Update(user *domain.User) error {
	res := u.db.Save(user)
	return res.Error
}

func (u UserRepository) Get(id uint64) *domain.User {
	var x UserEntity
	res := u.db.First(&x, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return x.ToUser()
	}
}

func (u UserRepository) GetByUuid(uid string) *domain.User {
	var tmp UserEntity
	res := u.db.Where("uuid = ?", uid).First(&tmp)
	if res.Error != nil {
		return nil
	}
	return tmp.ToUser()
}

func (u UserRepository) GetByUsername(username string) *domain.User {
	var tmp UserEntity
	res := u.db.Where("username = ?", username).First(&tmp)
	if res.Error != nil {
		return nil
	}
	return tmp.ToUser()
}

type GroupRepository struct {
	db *gorm.DB
}

func (g GroupRepository) Create(group *domain.Group) error {
	model := FromGroup(group)
	res := g.db.Create(model)
	if res.Error != nil {
		return res.Error
	}
	group.ID = model.ID
	group.CreatedAt = model.CreatedAt
	group.UpdatedAt = model.UpdatedAt
	return nil
}

func (g GroupRepository) Update(group *domain.Group) error {
	//TODO implement me
	panic("implement me")
}

func (g GroupRepository) Get(id uint64) *domain.Group {
	var x GroupEntity
	res := g.db.First(&x, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return x.ToGroup()
	}
}

func (g GroupRepository) UserGroupsPaginate(user uint64, page uint) Paginate[domain.Group] {
	var results []*domain.Group

	return Paginate[domain.Group]{
		Results:  results,
		Page:     page,
		NextPage: false,
	}
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}
