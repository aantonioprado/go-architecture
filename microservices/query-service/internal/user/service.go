package user

type ReadService struct {
	repo Repository
}

func NewReadService(repo Repository) *ReadService {
	return &ReadService{repo: repo}
}

func (s *ReadService) ListUsers() ([]*User, error) {
	return s.repo.FindAll()
}

func (s *ReadService) GetUser(id string) (*User, error) {
	return s.repo.FindByID(id)
}

func (s *ReadService) Replicate(user User) error {
	return s.repo.Upsert(&user)
}

func (s *ReadService) RemoveReplica(id string) error {
	return s.repo.Delete(id)
}
