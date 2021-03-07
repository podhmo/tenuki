package tenuki

import (
	"io"
	"net/http"
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
		T:      t,
		Config: DefaultConfig(options...),
	}
	if f.Client == nil {
		f.Client = &http.Client{Timeout: 1 * time.Second}
	}
	return f
}
func WithoutCapture() func(*Config) {
	return func(c *Config) {
		c.captureEnabled = false
	}
}
func WithWriteFile(basedir string) func(*Config) {
	return func(c *Config) {
		c.writeFileBaseDir = basedir
	}
}
func WithLayout(layout *capture.Layout) func(*Config) {
	return func(c *Config) {
		c.layout = layout
	}
}

func (f *Facade) NewRequest(
	method, url string, body io.Reader,
) *http.Request {
	t := f.T
	t.Helper()
	if url == "" {
		url = "http://example.net"
	}
	return NewRequest(t, method, url, body)
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
		ct := &CapturedTransport{
			T:                 t,
			CapturedTransport: f.NewCaptureTransport(t.Name()),
		}
		ct.CapturedTransport.Printer = ct
		ct.Transport = client.Transport
		client.Transport = ct
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

func NewRequest(
	t *testing.T,
	method, url string, body io.Reader,
) *http.Request {
	t.Helper()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("!! NewRequest: %+v", err)
	}
	return req
}
