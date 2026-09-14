package server

import (
	"net/http"

	"github.com/aantonioprado/go-architecture/layered/internal/handler"
	"github.com/aantonioprado/go-architecture/layered/internal/routes"
)

func Build() http.Handler {
	handlers := routes.Handlers{
		Health: handler.NewHealthHandler(),
	}

	return routes.NewRouter(handlers)
}
