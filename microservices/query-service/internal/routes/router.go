package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aantonioprado/go-architecture/microservices/query-service/internal/health"
	"github.com/aantonioprado/go-architecture/microservices/query-service/internal/middleware"
	"github.com/aantonioprado/go-architecture/microservices/query-service/internal/user"
)

type Handlers struct {
	Health   *health.Handler
	Users    *user.UserController
	Internal *user.InternalController
}

func NewRouter(h Handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.Users.ListUsers)
		r.Get("/{id}", h.Users.GetUser)
	})

	r.Route("/internal/users", func(r chi.Router) {
		r.Post("/", h.Internal.Replicate)
		r.Put("/{id}", h.Internal.UpdateReplica)
		r.Delete("/{id}", h.Internal.RemoveReplica)
	})

	r.Get("/health", h.Health.GetHealthCheck)

	return r
}
