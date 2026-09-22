package routes

import (
	"net/http"
	"net/http/httputil"

	"github.com/go-chi/chi/v5"

	"github.com/aantonioprado/go-architecture/microservices/gateway/internal/health"
	"github.com/aantonioprado/go-architecture/microservices/gateway/internal/middleware"
)

type Handlers struct {
	Health       *health.Handler
	CommandProxy *httputil.ReverseProxy
	QueryProxy   *httputil.ReverseProxy
}

func NewRouter(h Handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.CommandProxy.ServeHTTP)
		r.Put("/{id}", h.CommandProxy.ServeHTTP)
		r.Delete("/{id}", h.CommandProxy.ServeHTTP)
		r.Get("/", h.QueryProxy.ServeHTTP)
		r.Get("/{id}", h.QueryProxy.ServeHTTP)
	})

	r.Get("/health", h.Health.GetHealthCheck)

	return r
}
