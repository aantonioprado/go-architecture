package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/aantonioprado/go-architecture/microservices/query-service/internal/response"
)

type InternalController struct {
	service *ReadService
}

func NewInternalController(service *ReadService) *InternalController {
	return &InternalController{service: service}
}

func (c *InternalController) Replicate(w http.ResponseWriter, r *http.Request) {
	c.replicate(w, r, "")
}

func (c *InternalController) UpdateReplica(w http.ResponseWriter, r *http.Request) {
	c.replicate(w, r, chi.URLParam(r, "id"))
}

func (c *InternalController) replicate(w http.ResponseWriter, r *http.Request, id string) {
	var req ReplicaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}

	if id != "" {
		req.ID = id
	}

	if err := c.service.Replicate(toUser(req)); err != nil {
		writeError(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (c *InternalController) RemoveReplica(w http.ResponseWriter, r *http.Request) {
	if err := c.service.RemoveReplica(chi.URLParam(r, "id")); err != nil {
		writeError(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func toUser(req ReplicaRequest) User {
	createdAt, _ := time.Parse(time.RFC3339, req.CreatedAt)

	return User{
		ID:        req.ID,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: createdAt,
	}
}
