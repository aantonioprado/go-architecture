package repository

import "go-architecture-mvc/internal/model"

type UserRepository interface {
	CreateUser(user *model.User) error
	GetByIDUser(id string) (*model.User, error)
	GetAllUser() ([]*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id string) error
}
