package capture

import (
	"net/http"
)

type CapturedTransport struct {
	Transport http.RoundTripper
	Printer   printer

	DumpRequest  func(printer, *http.Request) error
	DumpResponse func(printer, *http.Response) error
}

func (ct *CapturedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := ct.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	dumpRequest := ct.DumpRequest
	if dumpRequest == nil {
		dumpRequest = DefaultDumper.DumpRequest
	}
	dumpResponse := ct.DumpResponse
	if dumpResponse == nil {
		dumpResponse = DefaultDumper.DumpResponse
	}

	if err := dumpRequest(ct.Printer, req); err != nil {
		return nil, err
	}
	res, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if err := dumpResponse(ct.Printer, res); err != nil {
		return nil, err
	}
	return res, nil
}

type printer interface {
	Printf(fmt string, args ...interface{})
}

type Dumper interface {
	DumpRequest(p printer, req *http.Request) error
	DumpResponse(p printer, res *http.Response) error
}

var DefaultDumper Dumper = &ConsoleDumper{}
