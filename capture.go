package tenuki

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

type CapturedTransport struct {
	Transport http.RoundTripper
	T         hasLogf
}

func (ct *CapturedTransport) Capture(t hasLogf) func() {
	ct.T = t
	return func() {
		ct.T = nil
	}
}

func (ct *CapturedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if ct.T == nil {
		fmt.Fprintln(os.Stderr, "!! CapturedTransport.T is not found !!")
		fmt.Fprintln(os.Stderr, "please use `defer transport.Capture(t)()`")
	}

	transport := ct.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	b, err := httputil.DumpRequest(req, true /* body */)
	if err != nil {
		return nil, err
	}

	// TODO: use Printf? and then this prefix is used in adapter
	ct.T.Logf("\x1b[5G\x1b[0K\x1b[90mrequest:\n%s\x1b[0m", string(b))

	res, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	b, err = httputil.DumpResponse(res, true /* body */)
	if err != nil {
		return nil, err
	}
	ct.T.Logf("\x1b[5G\x1b[0K\x1b[90mresponse:\n%s\x1b[0m", string(b))
	return res, nil
}

func ToLogf(p printer) hasLogf {
	return &logfAdapter{printer: p}
}

type printer interface {
	Printf(fmt string, args ...interface{})
}
type logfAdapter struct {
	printer printer
}

func (a *logfAdapter) Logf(fmt string, args ...interface{}) {
	a.printer.Printf(fmt, args...)
}

type hasLogf interface {
	Logf(fmt string, args ...interface{})
}
