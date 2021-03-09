package tenuki

import (
	"net/http"
	"testing"
)

// NewCloseConnectionHandler is a handler that closes the connecion suddenly. this is one of the test utilities.
func NewCloseConnectionHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("close connection by remote (for %s)", t.Name())
		c.Close()
	}
}
