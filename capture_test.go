package tenuki_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/podhmo/tenuki"
)

func TestCapture(t *testing.T) {
	transport := tenuki.NewCaptureTransport(t)

	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `{"message": "hello"}`)
		},
	))
	defer ts.Close()

	client := &http.Client{Transport: transport}

	t.Run("request 1", func(t *testing.T) {
		defer transport.Capture(t)()

		client.Get(ts.URL)
	})

	t.Run("request 2", func(t *testing.T) {
		defer transport.Capture(t)()

		req, _ := http.NewRequest("POST", ts.URL, strings.NewReader(`{"me": "foo"}`))
		client.Do(req)
	})
}
