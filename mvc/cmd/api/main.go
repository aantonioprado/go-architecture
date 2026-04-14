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

	handlers := routes.Handlers{
		Health: healthController,
	}

	r := routes.NewRouter(handlers)

	log.Printf("Starting server on port %s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatal(err)
	}
}
