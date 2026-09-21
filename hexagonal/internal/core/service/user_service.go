package service

import (
	"github.com/aantonioprado/go-architecture/hexagonal/internal/core/domain"
	"github.com/aantonioprado/go-architecture/hexagonal/internal/core/ports"
)

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(name, email string) (domain.User, error) {
	if _, err := s.repo.FindByEmail(email); err == nil {
		return domain.User{}, domain.ErrEmailTaken
	}

	user, err := domain.NewUser(name, email)
	if err != nil {
		return domain.User{}, err
	}

	if err := s.repo.Create(user); err != nil {
		return domain.User{}, err
	}

	return *user, nil
}

func (s *userService) ListUsers() ([]domain.User, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	items := make([]domain.User, 0, len(users))
	for _, user := range users {
		items = append(items, *user)
	}

	return items, nil
}

func (s *userService) GetUser(id string) (domain.User, error) {
	user, err := s.repo.FindById(id)
	if err != nil {
		return domain.User{}, err
	}

	return *user, nil
}

func (s *userService) UpdateUser(id, name, email string) (domain.User, error) {
	existing, err := s.repo.FindById(id)
	if err != nil {
		return domain.User{}, err
	}

	updated, err := existing.Update(name, email)
	if err != nil {
		return domain.User{}, err
	}

	if updated.Email != existing.Email {
		if _, err := s.repo.FindByEmail(updated.Email); err == nil {
			return domain.User{}, domain.ErrEmailTaken
		}
	}

	if err := s.repo.Update(updated); err != nil {
		return domain.User{}, err
	}

	return *updated, nil
}

func (s *userService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}
