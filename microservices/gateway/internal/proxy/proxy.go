package proxy

import (
	"net/http/httputil"
	"net/url"
)

// New builds a reverse proxy to a backend service. The gateway never
// inspects or rewrites the response: status code and body pass through
// exactly as the backend sent them.
func New(targetBaseURL string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetBaseURL)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(target), nil
}
