package server

import (
	"fmt"
	"net/http"

	"github.com/DexterLB/skgt_api/backend"
)

// Server is an HTTP server which serves the API
type Server struct {
	backend *backend.Backend
}

// NewServer returns a new server using the specified backend instance
func NewServer(backend *backend.Backend) *Server {
	return &Server{
		backend: backend,
	}
}

// ServeHTTP implements the HTTP handler interface
func (b *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc("/info", b.info)

	mux.ServeHTTP(w, r)
}

func (b *Server) info(w http.ResponseWriter, r *http.Request) {
	message, err := b.backend.Info()

	if err != nil {
		http.Error(w, fmt.Sprintf("Info failed: %s", err), 500)
		return
	}

	fmt.Fprintf(w, "%s", message)
}
