package tenuki

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/podhmo/tenuki/capture"
)

var (
	CaptureEnabledDefault   bool   = true
	CaptureWriteFileBaseDir string = ""
	globalFileDumpCounter   int64  = 0
)

func init() {
	if ok, _ := strconv.ParseBool(os.Getenv("NOCAPTURE")); ok {
		log.Println("CAPTURE_DISABLED is true, so deactivate tenuki.capture function")
		CaptureEnabledDefault = false
	}
	if ok, _ := strconv.ParseBool(os.Getenv("CAPTURE_DISABLED")); ok {
		log.Println("CAPTURE_DISABLED is true, so deactivate tenuki.capture function")
		CaptureEnabledDefault = false
	}
	if filename := os.Getenv("CAPTURE_WRITEFILE"); filename != "" {
		log.Println("CAPTURE_WRITEFILE is set, so activate the function writing capture output to files")
		CaptureWriteFileBaseDir = filename
	}
}

type Facade struct {
	T      *testing.T
	Client *http.Client

	captureEnabled   bool
	writeFileBaseDir string

	extractor *ExtractFacade
	mu        sync.Mutex
}

func New(t *testing.T, options ...func(*Facade)) *Facade {
	f := &Facade{
		T:                t,
		captureEnabled:   CaptureEnabledDefault,
		writeFileBaseDir: CaptureWriteFileBaseDir,
	}
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
		ct := &CapturedTransport{T: f.T}
		if f.writeFileBaseDir != "" {
			ct.Dumper = &capture.FileDumper{
				BaseDir: capture.Dir(f.writeFileBaseDir),
				Counter: &globalFileDumpCounter,
			}
		}
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
