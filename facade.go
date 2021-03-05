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
	CaptureEnabledDefault bool = true
)

func init() {
	if ok, _ := strconv.ParseBool(os.Getenv("NOCAPTURE")); ok {
		CaptureEnabledDefault = false
	}
}

type Facade struct {
	T      *testing.T
	Client *http.Client

	captureEnabled bool

	extractor *ExtractFacade
	mu        sync.Mutex
}

func New(t *testing.T, options ...func(*Facade)) *Facade {
	f := &Facade{T: t, captureEnabled: CaptureEnabledDefault}
	for _, opt := range options {
		opt(f)
	}
	if f.Client == nil {
		f.Client = &http.Client{Timeout: 1 * time.Second}
	}
	return f
}

var noop = func() {}

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

	client := f.Client
	if client == http.DefaultClient {
		panic("!! invalid: http.DefaultClient is used")
	}

	// TODO: not goroutine safe
	originalTransport := client.Transport
	if f.captureEnabled {
		transport := &CapturedTransport{T: f.T}
		transport.Transport = client.Transport
		client.Transport = transport
	}
	defer func() {
		f.Client.Transport = originalTransport
	}()

	res, err := client.Do(req)
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
