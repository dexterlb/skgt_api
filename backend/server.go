package backend

import (
	"fmt"
	"net/http"
)

type BackendServer struct {
	backend *Backend
}

func NewBackendServer(backend *Backend) *BackendServer {
	return &BackendServer{
		backend: backend,
	}
}

func (b *BackendServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc("/info", b.info)

	mux.ServeHTTP(w, r)
}

func (b *BackendServer) info(w http.ResponseWriter, r *http.Request) {
	message, err := b.backend.Info("42")

	if err != nil {
		http.Error(w, fmt.Sprintf("Info failed: %s", err), 500)
		return
	}

	fmt.Fprintf(w, "%s", message)
}
