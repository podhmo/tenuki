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

	Capture bool
	wrapped bool

	extractor *ExtractFacade
	mu        sync.Mutex
}

func New(t *testing.T) *Facade {
	return &Facade{T: t, Capture: captureDefault}
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
