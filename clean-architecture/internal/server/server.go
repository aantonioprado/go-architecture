package server

import (
	"net/http"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/adapters/controller"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/adapters/gateway"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/routes"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/usecases"
)

func Build() http.Handler {
	userRepository := gateway.NewInMemoryUserRepository()
	userInteractor := usecases.NewUserInteractor(userRepository)

	handlers := routes.Handlers{
		Health: controller.NewHealthController(),
		Users:  controller.NewUserController(userInteractor),
	}

	return routes.NewRouter(handlers)
}
