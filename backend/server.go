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

	mux.HandleFunc("/foo", b.foo)
	mux.HandleFunc("/get_age", b.getAge)

	mux.ServeHTTP(w, r)
}

func (b *BackendServer) foo(w http.ResponseWriter, r *http.Request) {
	message, err := b.backend.Foo()

	if err != nil {
		http.Error(w, fmt.Sprintf("Foo failed: %s", err), 500)
		return
	}

	fmt.Fprintf(w, "%s", message)
}

func (b *BackendServer) getAge(w http.ResponseWriter, r *http.Request) {
	person := r.URL.Query().Get("person")

	age, err := b.backend.GetAge(person)

	if err != nil {
		http.Error(w, fmt.Sprintf("GetAge failed: %s", err), 500)
		return
	}

	fmt.Fprintf(w, "%s's age is %d", person, age)
}
