package server

import (
	"net/http"

	"github.com/aantonioprado/go-architecture/layered/internal/handler"
	"github.com/aantonioprado/go-architecture/layered/internal/repository"
	"github.com/aantonioprado/go-architecture/layered/internal/routes"
	"github.com/aantonioprado/go-architecture/layered/internal/service"
)

func Build() http.Handler {
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository)

	handlers := routes.Handlers{
		Health: handler.NewHealthHandler(),
		Users:  handler.NewUserHandler(userService),
	}

	return routes.NewRouter(handlers)
}
