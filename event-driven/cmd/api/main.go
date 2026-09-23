package main

import (
	"log"
	"net/http"

	"github.com/aantonioprado/go-architecture/event-driven/internal/shared/config"
	"github.com/aantonioprado/go-architecture/event-driven/internal/shared/server"
)

func main() {
	cfg := config.Load()
	app := server.Build()

	log.Printf("Starting server on port %s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, app); err != nil {
		log.Fatal(err)
	}
}
