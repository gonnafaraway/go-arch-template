package http

import (
	"net/http"
	"time"
)

func NewHTTPTransport(port string) *http.Server {
	httpServer := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return httpServer
}
