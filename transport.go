package tenuki

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type RoundTripFunc func(*http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func HandlerTripperFunc(handle func(w http.ResponseWriter, r *http.Request)) *HandlerTripper {
	return &HandlerTripper{
		Handler: http.HandlerFunc(handle),
	}
}

type HandlerTripper struct {
	Before  func(*http.Request)
	Handler http.Handler
	After   func(*http.Response, *http.Request)
}

func (t *HandlerTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.Before != nil {
		t.Before(r)
	}
	w := httptest.NewRecorder()
	t.Handler.ServeHTTP(w, r)
	res := w.Result()
	if t.After != nil {
		t.After(res, r)
	}
	return res, nil
}

// NewErrorTransport returns a transport that returns error. this is one of the test utilities.
func NewErrorTransport(t *testing.T, genErr func() error) RoundTripFunc {
	return func(*http.Request) (*http.Response, error) {
		t.Helper()
		err := genErr()
		t.Logf("test helper -- returns error %T in transport (for %s) ..", err, t.Name())
		return nil, err
	}
}
