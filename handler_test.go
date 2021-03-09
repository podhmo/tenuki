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

func TestTimeoutHandler(t *testing.T) {
	handler, client := tenuki.NewTimeoutHandlerWithClient(t)
	ts := httptest.NewServer(handler)
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("unexpected request creation error %#+v", err)
	}

	_, err = client.Do(req)
	if err == nil {
		t.Fatal("must be error, but nil")
	}

	// &url.Error{Op:"Get", URL:"http://127.0.0.1:59233", Err:(*http.httpError)(0xc0001b2040)}
	for {
		inner, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}
		err = inner.Unwrap()
	}
	if err, ok := err.(interface{ Timeout() bool }); !(ok && err.Timeout()) {
		t.Errorf("timeout is expected but return error is %[1]T, %#+[1]v", err)
	}
}
