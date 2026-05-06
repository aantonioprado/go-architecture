package service

import (
	"go-architecture-mvc/internal/model"
	"go-architecture-mvc/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{
		repo: r,
	}
}

func (s *UserService) Create(name, email string) error {
	return nil
}

func (s *UserService) FindAll() ([]*model.User, error) {
	return nil, nil
}

func (s *UserService) FindByID(id string) (*model.User, error) {
	return nil, nil
}

func (s *UserService) Update(id, name, email string) error {
	return nil
}

func (s *UserService) Delete(id string) error {
	return nil
}
