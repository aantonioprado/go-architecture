package server

import (
	"net/http"

	"github.com/aantonioprado/go-architecture/microservices/command-service/internal/health"
	"github.com/aantonioprado/go-architecture/microservices/command-service/internal/replication"
	"github.com/aantonioprado/go-architecture/microservices/command-service/internal/routes"
	"github.com/aantonioprado/go-architecture/microservices/command-service/internal/user"
)

func Build(queryServiceURL string) (http.Handler, error) {
	replicator, err := replication.NewHTTPReplicator(queryServiceURL)
	if err != nil {
		return nil, err
	}

	userRepository := user.NewInMemoryUserRepository()
	userService := user.NewUserService(userRepository, replicator)
	userController := user.NewUserController(userService)

	handlers := routes.Handlers{
		Health: health.NewHandler(),
		Users:  userController,
	}

	return routes.NewRouter(handlers), nil
}
