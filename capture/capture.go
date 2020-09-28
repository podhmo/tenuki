package capture

import (
	"net/http"
	"net/http/httputil"
)

type CapturedTransport struct {
	Transport http.RoundTripper
	Printer   printer
}

type printer interface {
	Printf(fmt string, args ...interface{})
}

func (ct *CapturedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := ct.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	b, err := httputil.DumpRequest(req, true /* body */)
	if err != nil {
		return nil, err
	}

	ct.Printer.Printf("\x1b[90mrequest:\n%s\x1b[0m", string(b))

	res, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	b, err = httputil.DumpResponse(res, true /* body */)
	if err != nil {
		return nil, err
	}
	ct.Printer.Printf("\x1b[90mresponse:\n%s\x1b[0m", string(b))
	return res, nil
}
