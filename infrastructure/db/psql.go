package db

import (
	"errors"
	"financial/domain"
	"gorm.io/gorm"
	"time"
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
	model := &GroupEntity{
		InviteLink: group.InviteLink,
		CreatorId:  group.CreatorId,
		Name:       group.Name,
	}
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
	res := g.db.Model(&GroupEntity{}).Where("id = ?", group.ID).Updates(GroupEntity{InviteLink: group.InviteLink, Name: group.Name})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (g GroupRepository) Get(id uint64) *domain.Group {
	var x GroupEntity
	res := g.db.Preload("Creator").First(&x, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return x.ToGroup()
	}
}

func (g GroupRepository) UserGroupsPaginate(user uint64, page uint) *Paginate[domain.Group] {
	var results []*GroupEntity
	perPage := 2

	g.db.Where("creator_id = ?", user).
		Offset((int(page) - 1) * perPage).
		Limit(perPage + 1).
		Find(&results)

	groups := make([]*domain.Group, len(results))

	for i, result := range results {
		groups[i] = &domain.Group{
			InviteLink: result.InviteLink,
			DeletedAt:  result.DeletedAt,
			CreatedAt:  result.CreatedAt,
			UpdatedAt:  result.UpdatedAt,
			CreatorId:  result.CreatorId,
			Name:       result.Name,
			ID:         result.ID,
		}
	}

	return MakeSimplePaginate[domain.Group](groups, int(page), perPage)
}

func (g GroupRepository) SoftDelete(id uint64) error {
	res := g.db.Model(&GroupEntity{}).Where("id = ?", id).Update("deleted_at", time.Now())
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (g GroupRepository) UserLeaveGroup(groupId uint64, userId uint64) error {
	//TODO leave group (GroupMemberEntity need)
	res := g.db.Model(&GroupEntity{}).Where("id = ?", groupId).Update("deleted_at", userId)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}
