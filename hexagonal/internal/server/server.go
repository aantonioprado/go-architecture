package server

import (
	"net/http"

	httpadapter "github.com/aantonioprado/go-architecture/hexagonal/internal/adapters/primary/http"
	"github.com/aantonioprado/go-architecture/hexagonal/internal/adapters/secondary/memory"
	"github.com/aantonioprado/go-architecture/hexagonal/internal/core/service"
	"github.com/aantonioprado/go-architecture/hexagonal/internal/routes"
)

// Build is the composition root: it wires the secondary adapter into the core,
// and the core into the primary adapter, without either side knowing about the other.
func Build() http.Handler {
	userRepository := memory.NewInMemoryUserRepository()
	userService := service.NewUserService(userRepository)

	handlers := routes.Handlers{
		Health: httpadapter.NewHealthHandler(),
		Users:  httpadapter.NewUserHandler(userService),
	}

	return routes.NewRouter(handlers)
}
