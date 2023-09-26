package application

import (
	"errors"
	"financial/domain"
	"financial/infrastructure/db"
)

type IUserService interface {
	AddUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	GetUserById(id uint64) (*domain.User, error)
	GetUserByUuid(uuid string) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
}

type UserService struct {
	reader db.UserReader
	writer db.UserWriter
}

func NewUserService(reader db.UserReader, writer db.UserWriter) IUserService {
	return &UserService{reader: reader, writer: writer}
}

func (service UserService) AddUser(user *domain.User) error {
	err := service.writer.Create(user)
	if err != nil {
		return err
	}
	return nil
}

func (service UserService) UpdateUser(user *domain.User) error {
	err := service.writer.Update(user)
	if err != nil {
		return err
	}
	return nil
}

func (service UserService) GetUserById(id uint64) (*domain.User, error) {
	tmp := service.reader.Get(id)
	if tmp != nil {
		return tmp, nil
	}

	return nil, errors.New("Not found")
}

func (service UserService) GetUserByUuid(uuid string) (*domain.User, error) {
	tmp := service.reader.GetByUuid(uuid)
	if tmp != nil {
		return tmp, nil
	}

	return nil, errors.New("Not found")
}

func (service UserService) GetUserByUsername(username string) (*domain.User, error) {
	tmp := service.reader.GetByUsername(username)
	if tmp != nil {
		return tmp, nil
	}

	return nil, errors.New("Not found")
}
