package main

import (
	"log"
	"net/http"

	"github.com/aantonioprado/go-architecture/microservices/gateway/internal/config"
	"github.com/aantonioprado/go-architecture/microservices/gateway/internal/server"
)

func main() {
	cfg := config.Load()

	app, err := server.Build(cfg.CommandServiceURL, cfg.QueryServiceURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting gateway on port %s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, app); err != nil {
		log.Fatal(err)
	}
}
