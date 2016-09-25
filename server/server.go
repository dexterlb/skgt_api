package server

import (
	"fmt"
	"net/http"

	"github.com/DexterLB/skgt_api/backend"
)

// Server is an HTTP server which serves the API
type Server struct {
	backend *backend.Backend

	mux http.Handler
}

// New returns a new server using the specified backend instance
func New(backend *backend.Backend) *Server {
	mux := http.NewServeMux()
	s := &Server{
		backend: backend,
		mux:     mux,
	}

	mux.HandleFunc("/info", s.info)

	return s
}

// ServeHTTP implements the HTTP handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := s.checkAPIKey(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to verify api key: %s", err), 500)
		return
	}

	s.mux.ServeHTTP(w, r)
}

func (s *Server) info(w http.ResponseWriter, r *http.Request) {
	message, err := s.backend.Info()

	if err != nil {
		http.Error(w, fmt.Sprintf("Info failed: %s", err), 500)
		return
	}

	fmt.Fprintf(w, "%s", message)
}

func (s *Server) checkAPIKey(r *http.Request) error {
	apiKey := r.URL.Query().Get("api_key")
	if apiKey == "" {
		return fmt.Errorf("API Key is empty")
	}

	return s.backend.CheckAPIKey(apiKey)
}
