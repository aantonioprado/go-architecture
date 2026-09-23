package user

type UserService struct {
	repo      Repository
	publisher EventPublisher
}

func NewUserService(repo Repository, publisher EventPublisher) *UserService {
	return &UserService{repo: repo, publisher: publisher}
}

func (s *UserService) CreateUser(name, email string) (*User, error) {
	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, ErrEmailTaken
	}

	user, err := NewUser(name, email)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	s.publisher.Publish(UserCreated{
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})

	return user, nil
}

func (s *UserService) ListUsers() ([]*User, error) {
	return s.repo.FindAll()
}

func (s *UserService) GetUser(id string) (*User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) UpdateUser(id, name, email string) (*User, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	updated, err := existing.Update(name, email)
	if err != nil {
		return nil, err
	}

	if updated.Email != existing.Email {
		if _, err := s.repo.FindByEmail(updated.Email); err == nil {
			return nil, ErrEmailTaken
		}
	}

	if err := s.repo.Update(updated); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *UserService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}
