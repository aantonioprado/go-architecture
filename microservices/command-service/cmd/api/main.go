package main

import (
	"log"
	"net/http"

	"github.com/aantonioprado/go-architecture/microservices/command-service/internal/config"
	"github.com/aantonioprado/go-architecture/microservices/command-service/internal/server"
)

func main() {
	cfg := config.Load()

	app, err := server.Build(cfg.QueryServiceURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting command-service on port %s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, app); err != nil {
		log.Fatal(err)
	}
}
