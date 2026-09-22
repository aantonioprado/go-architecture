package user

import (
	"log"
	"sync"
)

type UserService struct {
	repo       Repository
	replicator Replicator
	wg         sync.WaitGroup
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

	s.replicateAsync("create", user.ID, *user, s.replicator.ReplicateCreate)

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

	s.replicateAsync("update", updated.ID, *updated, s.replicator.ReplicateUpdate)

	return updated, nil
}

func (s *UserService) DeleteUser(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	s.wg.Add(1)

	go func(id string) {
		defer s.wg.Done()
		defer recoverReplication("delete", id)

		if err := s.replicator.ReplicateDelete(id); err != nil {
			logReplicationErr("delete", id, err)
		}
	}(id)

	return nil
}

func (s *UserService) Wait() {
	s.wg.Wait()
}

func (s *UserService) replicateAsync(action, id string, u User, fn func(User) error) {
	s.wg.Add(1)

	go func(action, id string, u User) {
		defer s.wg.Done()
		defer recoverReplication(action, id)

		if err := fn(u); err != nil {
			logReplicationErr(action, id, err)
		}
	}(action, id, u)
}

func recoverReplication(action, id string) {
	if r := recover(); r != nil {
		log.Printf("[replication] handler for %s on user %s panicked: %v", action, id, r)
	}
}

func logReplicationErr(action, id string, err error) {
	log.Printf("[replication] failed to replicate %s for user %s: %v", action, id, err)
}
