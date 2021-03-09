package tenuki_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/podhmo/tenuki"
)

func TestCloseConnectionHandler(t *testing.T) {
	ts := httptest.NewServer(tenuki.NewCloseConnectionHandler(t))
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("unexpected request creation error %+v", err)
	}

	_, err = (&http.Client{Timeout: 100 * time.Millisecond}).Do(req)
	if err == nil {
		t.Fatal("must be error, but nil")
	}
	if !errors.Is(err, io.EOF) {
		t.Errorf("EOF is expected but return error is %[1]T, %+[1]v", err)
	}
}
