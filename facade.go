package tenuki

import (
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/podhmo/tenuki/capture"
)

type Facade struct {
	T *testing.T
	*Config

	Client    *http.Client
	extractor *ExtractFacade
	mu        sync.Mutex
}

func New(t *testing.T, options ...func(*Config)) *Facade {
	f := &Facade{
		T: t,
		Config: DefaultConfig(append([]func(*Config){func(c *Config) {
			c.disableCount = true
		}}, options...)...),
	}
	if f.Client == nil {
		f.Client = &http.Client{Timeout: 1 * time.Second}
	}
	return f
}

func (f *Facade) NewRequest(
	method, url string, body io.Reader,
) *http.Request {
	f.T.Helper()
	if !strings.Contains(url, "://") {
		url = "http://example.net/" + strings.TrimPrefix(url, "/")
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		f.T.Fatalf("!! NewRequest: %+v", err)
	}
	return req
}

func (f *Facade) NewJSONRequest(
	method, url string, body io.Reader,
) *http.Request {
	f.T.Helper()
	req := f.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json")
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
		ct := NewCaptureTransport(t, client.Transport)
		client.Transport = ct

		// for logf output (but this code is not good)
		if transport, ok := ct.Transport.(*capture.ConsoleTransport); ok {
			transport.Printer = ct
		}
	}
	defer func() {
		f.Client.Transport = originalTransport
	}()

	res, err := client.Do(req)
	if err != nil {
		if a.ExpectError == nil {
			t.Fatalf("!! Do: %+v", err)
		} else if err := a.ExpectError(err); err != nil {
			t.Fatalf("!! Do, ExpectError: %+v", err)
		}
	} else {
		if a.ExpectError != nil {
			t.Fatalf("!! Do, ExpectError: fatal is expected, but nil")
		}
	}
	for _, check := range a.Checks {
		check(t, res)
	}
	return res
}

func (f *Facade) DoHandler(
	handler http.Handler,
	req *http.Request,
	options ...AssertOption,
) *http.Response {
	client := f.Client
	if client == http.DefaultClient {
		panic("!! invalid: http.DefaultClient is used")
	}

	originalTransport := f.Client.Transport
	f.Client.Transport = &HandlerTripper{Handler: handler}
	res := f.Do(req, options...)
	f.Client.Transport = originalTransport
	return res
}
func (f *Facade) DoHandlerFunc(
	handler http.HandlerFunc,
	req *http.Request,
	options ...AssertOption,
) *http.Response {
	return f.DoHandler(handler, req, options...)
}

type Assertion struct {
	StatusCode  int
	Checks      []func(t *testing.T, res *http.Response)
	ExpectError func(err error) error
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
func AssertError(handler func(err error) error) AssertOption {
	return func(a *Assertion) {
		a.ExpectError = handler
	}
}

func NewAssertion() *Assertion {
	return &Assertion{}
}

type AssertOption func(*Assertion)
