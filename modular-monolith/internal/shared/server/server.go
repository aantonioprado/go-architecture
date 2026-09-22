package server

import (
	"net/http"

	"github.com/aantonioprado/go-architecture/modular-monolith/internal/health"
	"github.com/aantonioprado/go-architecture/modular-monolith/internal/shared/routes"
	"github.com/aantonioprado/go-architecture/modular-monolith/internal/user"
)

func Build() http.Handler {
	userRepository := user.NewInMemoryUserRepository()
	userService := user.NewUserService(userRepository)
	userController := user.NewUserController(userService)

	handlers := routes.Handlers{
		Health: health.NewHandler(),
		Users:  userController,
	}

	return routes.NewRouter(handlers)
}
