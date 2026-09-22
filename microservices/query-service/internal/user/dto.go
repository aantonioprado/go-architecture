package user

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

// ReplicaRequest is the wire shape command-service sends to /internal/users;
// unlike UserRequest elsewhere in this repository, it carries the full
// record (including ID and CreatedAt), since this service never generates
// those itself.
type ReplicaRequest struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}
