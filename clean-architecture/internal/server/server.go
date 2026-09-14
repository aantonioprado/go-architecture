package server

import (
	"net/http"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/interfaceadapters/controller"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/interfaceadapters/gateway"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/routes"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/usecases"
)

// Build is the composition root: it is the only place in this
// architecture that knows about a concrete gateway implementation.
func Build() http.Handler {
	userRepository := gateway.NewInMemoryUserRepository()
	userInteractor := usecases.NewUserInteractor(userRepository)

	handlers := routes.Handlers{
		Health: controller.NewHealthController(),
		Users:  controller.NewUserController(userInteractor),
	}

	return routes.NewRouter(handlers)
}
