package main

import (
	"log"
	"net/http"

	"github.com/aantonioprado/go-architecture/microservices/command-service/internal/config"
	"github.com/aantonioprado/go-architecture/microservices/command-service/internal/server"
)

func main() {
	cfg := config.Load()
	app := server.Build(cfg.QueryServiceURL)

	log.Printf("Starting command-service on port %s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, app); err != nil {
		log.Fatal(err)
	}
}
