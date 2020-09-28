package tenuki

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	captureDefault bool = true
)

func init() {
	if ok, _ := strconv.ParseBool(os.Getenv("NOCAPTURE")); ok {
		captureDefault = false
	}
}

type Facade struct {
	T      *testing.T
	Client *http.Client

	capture bool
	wrapped bool

	extractor *ExtractFacade
	mu        sync.Mutex
}

func New(t *testing.T) *Facade {
	return &Facade{T: t, capture: captureDefault}
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

	if f.capture {
		transport := &CapturedTransport{T: f.T}
		transport.Transport = client.Transport
		client.Transport = transport
	}
	f.Client = client
	return client
}

var noop = func() {}

func (f *Facade) Capture(t *testing.T) func() {
	if !f.capture {
		return noop
	}
	t.Helper()

	transport := f.client().Transport
	internal, ok := transport.(*CapturedTransport)
	if !ok {
		t.Fatalf("!! Capture: something wrong, transport is not captured")
		return noop
	}
	teardown := internal.Capture(t)
	return func() {
		f.mu.Lock()
		defer f.mu.Unlock()
		teardown()
		internal.T = f.T // rollback
	}
}

func (f *Facade) NewRequest(
	method, url string, body io.Reader,
) *http.Request {
	t := f.T
	t.Helper()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("!! NewRequest: %+v", err)
	}
	return req
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
