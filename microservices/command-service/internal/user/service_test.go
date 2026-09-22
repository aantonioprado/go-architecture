package user_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/aantonioprado/go-architecture/microservices/command-service/internal/user"
)

type fakeRepository struct {
	users map[string]*user.User
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{users: make(map[string]*user.User)}
}

func (r *fakeRepository) Create(u *user.User) error {
	r.users[u.ID] = u
	return nil
}

func (r *fakeRepository) FindAll() ([]*user.User, error) {
	users := make([]*user.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}

	return users, nil
}

func (r *fakeRepository) FindByID(id string) (*user.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, user.ErrUserNotFound
	}

	return u, nil
}

func (r *fakeRepository) Update(u *user.User) error {
	if _, ok := r.users[u.ID]; !ok {
		return user.ErrUserNotFound
	}

	r.users[u.ID] = u

	return nil
}

func (r *fakeRepository) Delete(id string) error {
	if _, ok := r.users[id]; !ok {
		return user.ErrUserNotFound
	}

	delete(r.users, id)

	return nil
}

func (r *fakeRepository) FindByEmail(email string) (*user.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, user.ErrUserNotFound
}

type fakeReplicator struct {
	mu      sync.Mutex
	created []user.User
	updated []user.User
	deleted []string
	failing bool
}

func (r *fakeReplicator) ReplicateCreate(u user.User) error {
	if r.failing {
		return errors.New("replication unavailable")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.created = append(r.created, u)

	return nil
}

func (r *fakeReplicator) ReplicateUpdate(u user.User) error {
	if r.failing {
		return errors.New("replication unavailable")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.updated = append(r.updated, u)

	return nil
}

func (r *fakeReplicator) ReplicateDelete(id string) error {
	if r.failing {
		return errors.New("replication unavailable")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.deleted = append(r.deleted, id)

	return nil
}

func TestUserService_CreateUser(t *testing.T) {
	rep := &fakeReplicator{}
	svc := user.NewUserService(newFakeRepository(), rep)

	u, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.Email != "antonio@antonioeprado.dev" {
		t.Errorf("expected email %q, got %q", "antonio@antonioeprado.dev", u.Email)
	}

	svc.Wait()

	if len(rep.created) != 1 {
		t.Fatalf("expected 1 replicated create, got %d", len(rep.created))
	}

	if rep.created[0].ID != u.ID {
		t.Errorf("expected replicated ID %q, got %q", u.ID, rep.created[0].ID)
	}
}

type blockingReplicator struct {
	release chan struct{}
}

func (r *blockingReplicator) ReplicateCreate(u user.User) error {
	<-r.release
	return nil
}

func (r *blockingReplicator) ReplicateUpdate(u user.User) error {
	<-r.release
	return nil
}

func (r *blockingReplicator) ReplicateDelete(id string) error {
	<-r.release
	return nil
}

func TestUserService_CreateUser_DoesNotWaitForReplication(t *testing.T) {
	rep := &blockingReplicator{release: make(chan struct{})}
	svc := user.NewUserService(newFakeRepository(), rep)

	done := make(chan struct{})
	go func() {
		if _, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev"); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("CreateUser blocked waiting for replication to finish")
	}

	close(rep.release)
	svc.Wait()
}

func TestUserService_CreateUser_ReplicationFailureDoesNotFailWrite(t *testing.T) {
	rep := &fakeReplicator{failing: true}
	svc := user.NewUserService(newFakeRepository(), rep)

	u, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("expected write to succeed even if replication fails, got: %v", err)
	}

	if u == nil {
		t.Fatal("expected a created user")
	}

	svc.Wait()
}

func TestUserService_CreateUser_MissingName(t *testing.T) {
	svc := user.NewUserService(newFakeRepository(), &fakeReplicator{})

	if _, err := svc.CreateUser("", "antonio@antonioeprado.dev"); !errors.Is(err, user.ErrNameRequired) {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	svc := user.NewUserService(newFakeRepository(), &fakeReplicator{})

	if _, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := svc.CreateUser("Outro", "antonio@antonioeprado.dev"); !errors.Is(err, user.ErrEmailTaken) {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	rep := &fakeReplicator{}
	svc := user.NewUserService(newFakeRepository(), rep)

	created, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svc.Wait()

	updated, err := svc.UpdateUser(created.ID, "Antônio Elias Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated.Name != "Antônio Elias Prado" {
		t.Errorf("expected updated name, got %q", updated.Name)
	}

	svc.Wait()

	if len(rep.updated) != 1 {
		t.Fatalf("expected 1 replicated update, got %d", len(rep.updated))
	}
}

func TestUserService_UpdateUser_NotFound(t *testing.T) {
	svc := user.NewUserService(newFakeRepository(), &fakeReplicator{})

	if _, err := svc.UpdateUser("unknown-id", "Antônio Prado", "antonio@antonioeprado.dev"); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_UpdateUser_DuplicateEmail(t *testing.T) {
	svc := user.NewUserService(newFakeRepository(), &fakeReplicator{})

	if _, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second, err := svc.CreateUser("Outro", "outro@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := svc.UpdateUser(second.ID, "Outro", "antonio@antonioeprado.dev"); !errors.Is(err, user.ErrEmailTaken) {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	rep := &fakeReplicator{}
	svc := user.NewUserService(newFakeRepository(), rep)

	created, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svc.Wait()

	if err := svc.DeleteUser(created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	svc.Wait()

	if len(rep.deleted) != 1 || rep.deleted[0] != created.ID {
		t.Fatalf("expected replicated delete for %q, got %v", created.ID, rep.deleted)
	}
}

func TestUserService_DeleteUser_NotFound(t *testing.T) {
	svc := user.NewUserService(newFakeRepository(), &fakeReplicator{})

	if err := svc.DeleteUser("unknown-id"); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
