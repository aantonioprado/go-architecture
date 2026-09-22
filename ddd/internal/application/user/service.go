package user

import (
	"log"
	"time"

	"github.com/aantonioprado/go-architecture/ddd/internal/domain/user"
)

type Service interface {
	CreateUser(name, email string) (*user.User, error)
	ListUsers() ([]*user.User, error)
	GetUser(id string) (*user.User, error)
	UpdateUser(id, name, email string) (*user.User, error)
	DeleteUser(id string) error
}

type applicationService struct {
	repo user.Repository
}

func NewService(repo user.Repository) Service {
	return &applicationService{repo: repo}
}

func (s *applicationService) CreateUser(name, email string) (*user.User, error) {
	emailVO, err := user.NewEmail(email)
	if err != nil {
		return nil, err
	}

	if _, err := s.repo.FindByEmail(emailVO); err == nil {
		return nil, user.ErrEmailTaken
	}

	u, err := user.Register(name, emailVO)
	if err != nil {
		return nil, err
	}

	return s.saveAndPublish(u)
}

func (s *applicationService) ListUsers() ([]*user.User, error) {
	return s.repo.FindAll()
}

func (s *applicationService) GetUser(id string) (*user.User, error) {
	return s.repo.FindByID(id)
}

func (s *applicationService) UpdateUser(id, name, email string) (*user.User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	emailVO, err := user.NewEmail(email)
	if err != nil {
		return nil, err
	}

	if emailVO != u.Email() {
		if _, err := s.repo.FindByEmail(emailVO); err == nil {
			return nil, user.ErrEmailTaken
		}
	}

	if err := u.ChangeDetails(name, emailVO); err != nil {
		return nil, err
	}

	return s.saveAndPublish(u)
}

func (s *applicationService) DeleteUser(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	logEvent(user.UserDeleted{UserID: id, At: time.Now()})

	return nil
}

func (s *applicationService) saveAndPublish(u *user.User) (*user.User, error) {
	if err := s.repo.Save(u); err != nil {
		return nil, err
	}

	for _, event := range u.Events() {
		logEvent(event)
	}

	u.ClearEvents()

	return u, nil
}

func logEvent(event user.Event) {
	log.Printf("[domain-event] %T: %+v", event, event)
}
