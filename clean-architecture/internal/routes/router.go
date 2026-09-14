package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/interfaceadapters/controller"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/middleware"
)

type Handlers struct {
	Health *controller.HealthController
	Users  *controller.UserController
}

func NewRouter(h Handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.Users.CreateUser)
	})

	r.Get("/health", h.Health.GetHealthCheck)

	return r
}
