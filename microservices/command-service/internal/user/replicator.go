package user

type Replicator interface {
	ReplicateCreate(user User) error
	ReplicateUpdate(user User) error
	ReplicateDelete(id string) error
}
