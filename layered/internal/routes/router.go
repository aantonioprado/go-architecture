package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aantonioprado/go-architecture/layered/internal/handler"
	"github.com/aantonioprado/go-architecture/layered/internal/middleware"
)

type Handlers struct {
	Health *handler.HealthHandler
}

func NewRouter(h Handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/health", h.Health.GetHealthCheck)

	return r
}
