package service

import "go-architecture-mvc/internal/repository"

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
