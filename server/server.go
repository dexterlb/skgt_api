package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DexterLB/htmlparsing"
	"github.com/DexterLB/skgt_api/backend"
	"github.com/julienschmidt/httprouter"
)

// Server is an HTTP server which serves the API
type Server struct {
	backend        *backend.Backend
	parserSettings *htmlparsing.Settings

	router *httprouter.Router
}

// New returns a new server using the specified backend instance
func New(backend *backend.Backend, parserSettings *htmlparsing.Settings) *Server {
	router := httprouter.New()
	s := &Server{
		backend:        backend,
		parserSettings: parserSettings,
		router:         router,
	}

	router.GET("/info", s.info)
	router.GET("/stop/:stop_id/arrivals/realtime", jsonHandler(s.realtimeArrivals))
	router.GET("/transport/list", jsonHandler(s.transports))

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
	var apiKey string

	user, password, hasAuth := r.BasicAuth()
	if hasAuth {
		if user != "" {
			apiKey = user
		}
		if password != "" {
			apiKey = password
		}
	} else {
		apiKey = r.Header.Get("X-Api-Key")
		if apiKey == "" {
			apiKey = r.URL.Query().Get("api_key")
		}
	}

	if apiKey == "" {
		return fmt.Errorf("API Key is empty")
	}

	return s.backend.CheckAPIKey(apiKey)
}

// jsonHandler wraps a function which returns JSON-marshable data (or an error)
// and returns a httrouter Handle which calls the function upon a request
func jsonHandler(handler func(params httprouter.Params) (interface{}, error)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		object, err := handler(params)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("error handling request: %s", err),
				http.StatusInternalServerError,
			)
		}

		data, err := json.MarshalIndent(object, "", "    ")
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("error marshaling data: %s", err),
				http.StatusInternalServerError,
			)
		}

		w.Write(data)
	}
}
