package server

import (
	"net/http"

	"github.com/aantonioprado/go-architecture/microservices/gateway/internal/health"
	"github.com/aantonioprado/go-architecture/microservices/gateway/internal/proxy"
	"github.com/aantonioprado/go-architecture/microservices/gateway/internal/routes"
)

func Build(commandServiceURL, queryServiceURL string) (http.Handler, error) {
	commandProxy, err := proxy.New(commandServiceURL)
	if err != nil {
		return nil, err
	}

	queryProxy, err := proxy.New(queryServiceURL)
	if err != nil {
		return nil, err
	}

	handlers := routes.Handlers{
		Health:       health.NewHandler(),
		CommandProxy: commandProxy,
		QueryProxy:   queryProxy,
	}

	return routes.NewRouter(handlers), nil
}
