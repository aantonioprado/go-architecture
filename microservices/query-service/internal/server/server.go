package server

import (
	"net/http"

	"github.com/aantonioprado/go-architecture/microservices/query-service/internal/health"
	"github.com/aantonioprado/go-architecture/microservices/query-service/internal/routes"
	"github.com/aantonioprado/go-architecture/microservices/query-service/internal/user"
)

func Build() http.Handler {
	userRepository := user.NewInMemoryUserRepository()
	readService := user.NewReadService(userRepository)

	handlers := routes.Handlers{
		Health:   health.NewHandler(),
		Users:    user.NewUserController(readService),
		Internal: user.NewInternalController(readService),
	}

	return routes.NewRouter(handlers)
}
