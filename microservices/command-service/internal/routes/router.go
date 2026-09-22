package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aantonioprado/go-architecture/microservices/command-service/internal/health"
	"github.com/aantonioprado/go-architecture/microservices/command-service/internal/middleware"
	"github.com/aantonioprado/go-architecture/microservices/command-service/internal/user"
)

type Handlers struct {
	Health *health.Handler
	Users  *user.UserController
}

func NewRouter(h Handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.Users.CreateUser)
		r.Put("/{id}", h.Users.UpdateUser)
		r.Delete("/{id}", h.Users.DeleteUser)
	})

	r.Get("/health", h.Health.GetHealthCheck)

	return r
}
