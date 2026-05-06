package repository

import "github.com/aantonioprado/go-architecture/mvc/internal/model"

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *model.User) error {
	return nil
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	return nil, nil
}

func (r *UserRepository) FindAll() ([]*model.User, error) {
	return nil, nil
}

func (r *UserRepository) Update(user *model.User) error {
	return nil
}

func (r *UserRepository) Delete(id string) error {
	return nil
}
