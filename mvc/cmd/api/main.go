package main

import (
	"go-architecture-mvc/internal/config"
	"go-architecture-mvc/internal/server"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()
	app := server.Build()

	log.Printf("Starting server on port %s", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, app); err != nil {
		log.Fatal(err)
	}
}
