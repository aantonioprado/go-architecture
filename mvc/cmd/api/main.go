package main

import (
	"go-architecture-mvc/internal/config"
	"go-architecture-mvc/internal/controller"
	"go-architecture-mvc/internal/routes"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	healthController := controller.NewHealthController()
	userController := controller.NewUserController()

	handlers := routes.Handlers{
		Health: healthController,
		Users:  userController,
	}

	r := routes.NewRouter(handlers)

	log.Printf("Starting server on port %s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
