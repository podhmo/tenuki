package tenuki_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/podhmo/tenuki"
)

func TestCapture(t *testing.T) {
	transport := tenuki.NewCaptureTransport(t, nil)

	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `{"message": "hello"}`)
		},
	))
	defer ts.Close()

	client := &http.Client{Transport: transport, Timeout: 1 * time.Second}

	t.Run("request 1", func(t *testing.T) {
		defer transport.Capture(t)()

		if _, err := client.Get(ts.URL); err != nil {
			t.Fatalf("!! %+v", err)
		}
	})

	t.Run("request 2", func(t *testing.T) {
		defer transport.Capture(t)()

		req, _ := http.NewRequest("POST", ts.URL, strings.NewReader(`{"me": "foo"}`))
		if _, err := client.Do(req); err != nil {
			t.Fatalf("!! %+v", err)
		}
	})
}
