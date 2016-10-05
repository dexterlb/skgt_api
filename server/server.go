package server

import (
	"fmt"
	"net/http"

	"github.com/DexterLB/skgt_api/backend"
	"github.com/julienschmidt/httprouter"
)

// Server is an HTTP server which serves the API
type Server struct {
	backend *backend.Backend

	router *httprouter.Router
}

// New returns a new server using the specified backend instance
func New(backend *backend.Backend) *Server {
	router := httprouter.New()
	s := &Server{
		backend: backend,
		router:  router,
	}

	router.GET("/info", s.info)

	return s
}

// ServeHTTP implements the HTTP handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := s.checkAPIKey(r)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("unable to verify api key: %s", err),
			http.StatusForbidden,
		)
		return
	}

	s.router.ServeHTTP(w, r)
}

func (s *Server) info(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
