package replication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aantonioprado/go-architecture/microservices/command-service/internal/user"
)

// HTTPReplicator keeps query-service's read model in sync over real HTTP
// calls, the only thing that connects the two services.
type HTTPReplicator struct {
	baseURL string
	client  *http.Client
}

func NewHTTPReplicator(baseURL string) *HTTPReplicator {
	return &HTTPReplicator{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

type replicaPayload struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

func (r *HTTPReplicator) ReplicateCreate(u user.User) error {
	return r.send(http.MethodPost, r.baseURL+"/internal/users", toReplicaPayload(u))
}

func (r *HTTPReplicator) ReplicateUpdate(u user.User) error {
	return r.send(http.MethodPut, r.baseURL+"/internal/users/"+u.ID, toReplicaPayload(u))
}

func (r *HTTPReplicator) ReplicateDelete(id string) error {
	return r.send(http.MethodDelete, r.baseURL+"/internal/users/"+id, nil)
}

func (r *HTTPReplicator) send(method, url string, payload any) error {
	var body *bytes.Buffer

	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(data)
	} else {
		body = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return fmt.Errorf("query-service returned status %d", res.StatusCode)
	}

	return nil
}

func toReplicaPayload(u user.User) replicaPayload {
	return replicaPayload{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}
