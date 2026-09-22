package user

import "log"

type UserService struct {
	repo       Repository
	replicator Replicator
}

func NewUserService(repo Repository, replicator Replicator) *UserService {
	return &UserService{repo: repo, replicator: replicator}
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

	if err := s.replicator.ReplicateCreate(*user); err != nil {
		logReplicationErr("create", user.ID, err)
	}

	return user, nil
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

	if err := s.replicator.ReplicateUpdate(*updated); err != nil {
		logReplicationErr("update", updated.ID, err)
	}

	return updated, nil
}

func (s *UserService) DeleteUser(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	if err := s.replicator.ReplicateDelete(id); err != nil {
		logReplicationErr("delete", id, err)
	}

	return nil
}

func logReplicationErr(action, id string, err error) {
	log.Printf("[replication] failed to replicate %s for user %s: %v", action, id, err)
}
