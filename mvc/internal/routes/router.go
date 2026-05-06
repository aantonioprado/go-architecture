package routes

import (
	"go-architecture-mvc/internal/controller"
	"go-architecture-mvc/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
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
		r.Get("/", h.Users.ListUsers)
		r.Get("/{id}", h.Users.FindUserByID)
	})

	r.Get("/health", h.Health.GetHealthCheck)

	return r
}
