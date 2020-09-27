package reqtest

import (
	"io"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
)

type Facade struct {
	T      *testing.T
	Client *http.Client

	Capture bool
	wrapped bool

	extractor *ExtractFacade
	mu        sync.Mutex
}

func New(t *testing.T) *Facade {
	return &Facade{T: t, Capture: true}
}

func (f *Facade) client() *http.Client {
	client := f.Client
	if client != nil {
		if f.wrapped {
			return client
		}
		if client == http.DefaultClient {
			panic("!! invalid: http.DefaultClient is used")
		}
	}

	f.wrapped = true

	if client == nil {
		client = &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   10 * time.Second, // xxx
		}
	}

	if f.Capture {
		client.Transport = &CapturedTransport{
			T:         f.T,
			Transport: client.Transport,
		}
	}
	f.Client = client
	return client
}

type RequestOption func(*http.Request)

func (f *Facade) NewRequest(
	method, url string, body io.Reader,
	options ...RequestOption,
) *http.Request {
	t := f.T
	t.Helper()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("!! NewRequest: %+v", err)
	}
	for _, opt := range options {
		opt(req)
	}
	return req
}

func WithResetQuery(modify func(u url.Values)) RequestOption {
	return func(req *http.Request) {
		var q url.Values
		modify(q)
		req.URL.RawQuery = q.Encode()
	}
}
func WithQuery(modify func(u url.Values)) RequestOption {
	return func(req *http.Request) {
		q := req.URL.Query()
		modify(q)
		req.URL.RawQuery = q.Encode()
	}
}

func (f *Facade) Do(
	req *http.Request,
	options ...AssertOption,
) *http.Response {
	t := f.T
	t.Helper()

	a := NewAssertion()
	for _, opt := range options {
		opt(a)
	}

	res, err := f.client().Do(req)
	if err != nil {
		t.Fatalf("!! Do: %+v", err)
	}

	for _, check := range a.Checks {
		check(t, res)
	}
	return res
}

type Assertion struct {
	StatusCode int
	Checks     []func(t *testing.T, res *http.Response)
}

func AssertStatus(code int) AssertOption {
	return func(a *Assertion) {
		a.StatusCode = code
		a.Checks = append(a.Checks, func(t *testing.T, res *http.Response) {
			t.Helper()
			if a.StatusCode != res.StatusCode {
				t.Errorf("status code:\nwant\n\t%+v\nbut\n\t%+v", a.StatusCode, res.StatusCode)
			}
		})
	}
}

func NewAssertion() *Assertion {
	return &Assertion{}
}

type AssertOption func(*Assertion)
