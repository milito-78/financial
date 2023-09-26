package psql

import (
	"errors"
	"financial/domain"
	"financial/infrastructure/db/psql/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u UserRepository) Create(user *domain.User) error {
	model := models.FromUser(user)
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
	var x models.UserEntity
	res := u.db.First(&x, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return x.ToUser()
	}
}

func (u UserRepository) GetByUuid(uid string) *domain.User {
	var tmp models.UserEntity
	res := u.db.Where("uuid = ?", uid).First(&tmp)
	if res.Error != nil {
		return nil
	}
	return tmp.ToUser()
}

func (u UserRepository) GetByUsername(username string) *domain.User {
	var tmp models.UserEntity
	res := u.db.Where("username = ?", username).First(&tmp)
	if res.Error != nil {
		return nil
	}
	return tmp.ToUser()
}
