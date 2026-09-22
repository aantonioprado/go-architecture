package server

import (
	"net/http"

	"github.com/aantonioprado/go-architecture/event-driven/internal/audit"
	"github.com/aantonioprado/go-architecture/event-driven/internal/health"
	"github.com/aantonioprado/go-architecture/event-driven/internal/notification"
	"github.com/aantonioprado/go-architecture/event-driven/internal/shared/events"
	"github.com/aantonioprado/go-architecture/event-driven/internal/shared/routes"
	"github.com/aantonioprado/go-architecture/event-driven/internal/user"
)

func Build() http.Handler {
	bus := events.NewBus()

	notification.NewListener(bus)
	audit.NewListener(bus)

	userRepository := user.NewInMemoryUserRepository()
	userService := user.NewUserService(userRepository, bus)
	userController := user.NewUserController(userService)

	handlers := routes.Handlers{
		Health: health.NewHandler(),
		Users:  userController,
	}

	return routes.NewRouter(handlers)
}
