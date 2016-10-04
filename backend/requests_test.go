package backend

import "testing"

func TestBackend_Transports(t *testing.T) {
	backend := fillDatabase(t)
	defer closeBackend(t, backend)
}
