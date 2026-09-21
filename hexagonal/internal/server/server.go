package server

import (
	"net/http"

	httpadapter "github.com/aantonioprado/go-architecture/hexagonal/internal/adapters/primary/http"
	"github.com/aantonioprado/go-architecture/hexagonal/internal/adapters/secondary/memory"
	"github.com/aantonioprado/go-architecture/hexagonal/internal/core/service"
	"github.com/aantonioprado/go-architecture/hexagonal/internal/routes"
)

func Build() http.Handler {
	userRepository := memory.NewInMemoryUserRepository()
	userService := service.NewUserService(userRepository)

	handlers := routes.Handlers{
		Health: httpadapter.NewHealthHandler(),
		Users:  httpadapter.NewUserHandler(userService),
	}

	return routes.NewRouter(handlers)
}
