package routes

import (
	"go-architecture-mvc/internal/controller"
	"go-architecture-mvc/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Health *controller.HealthController
}

func NewRouter(h Handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/health", h.Health.GetHealthCheck)

	return r
}
