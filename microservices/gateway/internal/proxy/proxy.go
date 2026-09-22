package proxy

import (
	"net/http/httputil"
	"net/url"
)

func New(targetBaseURL string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetBaseURL)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(target), nil
}
