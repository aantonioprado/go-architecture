package replication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/aantonioprado/go-architecture/microservices/command-service/internal/user"
)

type HTTPReplicator struct {
	baseURL string
	client  *http.Client
}

func NewHTTPReplicator(baseURL string) (*HTTPReplicator, error) {
	if _, err := url.Parse(baseURL); err != nil {
		return nil, err
	}

	return &HTTPReplicator{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}, nil
}

func (r *HTTPReplicator) ReplicateCreate(u user.User) error {
	return r.send(http.MethodPost, r.baseURL+"/internal/users", user.ToUserResponse(&u))
}

func (r *HTTPReplicator) ReplicateUpdate(u user.User) error {
	return r.send(http.MethodPut, r.baseURL+"/internal/users/"+u.ID, user.ToUserResponse(&u))
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
