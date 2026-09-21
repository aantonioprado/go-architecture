package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	httpadapter "github.com/aantonioprado/go-architecture/hexagonal/internal/adapters/primary/http"
	"github.com/aantonioprado/go-architecture/hexagonal/internal/middleware"
)

type Handlers struct {
	Health *httpadapter.HealthHandler
	Users  *httpadapter.UserHandler
}

func NewRouter(h Handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.Users.CreateUser)
		r.Get("/", h.Users.ListUsers)
		r.Get("/{id}", h.Users.GetUser)
		r.Put("/{id}", h.Users.UpdateUser)
		r.Delete("/{id}", h.Users.DeleteUser)
	})

	r.Get("/health", h.Health.GetHealthCheck)

	return r
}
