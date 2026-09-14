package server

import (
	"net/http"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/interfaceadapters/controller"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/routes"
)

func Build() http.Handler {
	handlers := routes.Handlers{
		Health: controller.NewHealthController(),
	}

	return routes.NewRouter(handlers)
}
