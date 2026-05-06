package server

import (
	"github.com/aantonioprado/go-architecture/mvc/internal/controller"
	"github.com/aantonioprado/go-architecture/mvc/internal/repository"
	"github.com/aantonioprado/go-architecture/mvc/internal/routes"
	"net/http"
)

func Build() http.Handler {
	userRepository := repository.NewUserRepository()

	handlers := routes.Handlers{
		Health: controller.NewHealthController(),
		Users:  controller.NewUserController(userRepository),
	}

	return routes.NewRouter(handlers)
}
