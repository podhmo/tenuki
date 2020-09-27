package reqtest

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"testing"
)

type CapturedTransport struct {
	Transport http.RoundTripper
	T         *testing.T
}

func (ct *CapturedTransport) Capture(t *testing.T) func() {
	ct.T = t
	return func() {
		ct.T = nil
	}
}

func (ct *CapturedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if ct.T == nil {
		fmt.Fprintln(os.Stderr, "!! CapturedTransport.T is not found. please use !!")
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
