package tenuki

import (
	"net/http"
	"testing"
	"time"
)

// NewCloseConnectionHandler is a handler that closes the connecion suddenly. this is one of the test utilities.
func NewCloseConnectionHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("test helper -- the connection is closed suddenly by remote side (for %s)", t.Name())
		c.Close()
	}
}

// NewTimeoutHandler returns a handler that closes the connecion suddenly. this is one of the test utilities.
func NewTimeoutHandler(t *testing.T, timeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Logf("test helper -- the timeout %+v is reached in remote side (for %s) ..", timeout, t.Name())
		time.Sleep(2 * timeout)
	}
}

// NewTimeoutHandlerWithClient returns a handler that timeout is occured. this is one of the test utilities.
func NewTimeoutHandlerWithClient(t *testing.T, timeouts ...time.Duration) (http.HandlerFunc, *http.Client) {
	t.Helper()
	timeout := 100 * time.Millisecond
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}
	return NewTimeoutHandler(t, timeout), &http.Client{Timeout: timeout}
}
