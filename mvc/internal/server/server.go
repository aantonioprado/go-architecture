package server

import (
	"go-architecture-mvc/internal/controller"
	"go-architecture-mvc/internal/repository"
	"go-architecture-mvc/internal/routes"
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
