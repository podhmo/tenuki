package tenuki

import (
	"net/http"
	"net/http/httptest"
)

type RoundTripFunc func(*http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
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
