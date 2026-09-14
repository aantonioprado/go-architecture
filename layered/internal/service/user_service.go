package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/aantonioprado/go-architecture/layered/internal/model"
	"github.com/aantonioprado/go-architecture/layered/internal/repository"
)

var (
	ErrNameRequired  = errors.New("name is required")
	ErrEmailRequired = errors.New("email is required")
	ErrEmailTaken    = errors.New("email already in use")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Create(name, email string) (*model.User, error) {
	if err := validate(name, email); err != nil {
		return nil, err
	}

	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, ErrEmailTaken
	}

	user := &model.User{
		ID:        newID(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) List() ([]*model.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) GetById(id string) (*model.User, error) {
	return s.repo.FindById(id)
}

func (s *UserService) Update(id, name, email string) (*model.User, error) {
	existing, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	if err := validate(name, email); err != nil {
		return nil, err
	}

	if email != existing.Email {
		if _, err := s.repo.FindByEmail(email); err == nil {
			return nil, ErrEmailTaken
		}
	}

	updated := &model.User{
		ID:        existing.ID,
		Name:      name,
		Email:     email,
		CreatedAt: existing.CreatedAt,
	}

	if err := s.repo.Update(updated); err != nil {
		return nil, err
	}

	return updated, nil
}

func validate(name, email string) error {
	if name == "" {
		return ErrNameRequired
	}

	if email == "" {
		return ErrEmailRequired
	}

	return nil
}

func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
