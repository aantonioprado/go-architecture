package server

import (
	"net/http"

	appuser "github.com/aantonioprado/go-architecture/ddd/internal/application/user"
	"github.com/aantonioprado/go-architecture/ddd/internal/infrastructure/persistence/memory"
	httpinterface "github.com/aantonioprado/go-architecture/ddd/internal/interfaces/http"
	"github.com/aantonioprado/go-architecture/ddd/internal/routes"
)

func Build() http.Handler {
	userRepository := memory.NewUserRepository()
	userService := appuser.NewService(userRepository)

	handlers := routes.Handlers{
		Health: httpinterface.NewHealthHandler(),
		Users:  httpinterface.NewUserHandler(userService),
	}

	return routes.NewRouter(handlers)
}
