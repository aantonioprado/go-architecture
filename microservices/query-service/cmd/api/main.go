package main

import (
	"log"
	"net/http"

	"github.com/aantonioprado/go-architecture/microservices/query-service/internal/config"
	"github.com/aantonioprado/go-architecture/microservices/query-service/internal/server"
)

func main() {
	cfg := config.Load()
	app := server.Build()

	log.Printf("Starting query-service on port %s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, app); err != nil {
		log.Fatal(err)
	}
}
