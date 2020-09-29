package capture

import (
	"net/http"
)

type CapturedTransport struct {
	Transport http.RoundTripper
	Printer   printer

	Dumper Dumper
}

func (ct *CapturedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := ct.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	dumper := ct.Dumper
	if dumper == nil {
		dumper = DefaultDumper
	}

	s, err := dumper.DumpRequest(ct.Printer, req)
	if err != nil {
		return nil, err
	}
	res, err := transport.RoundTrip(req)
	if err != nil {
		return nil, dumper.DumpError(ct.Printer, s, err)
	}
	if err := dumper.DumpResponse(ct.Printer, s, res); err != nil {
		return nil, err
	}
	return res, nil
}

type printer interface {
	Printf(fmt string, args ...interface{})
}

type State interface {
	Request() *http.Request
}
type Dumper interface {
	DumpRequest(p printer, req *http.Request) (State, error)
	DumpResponse(p printer, state State, res *http.Response) error
	DumpError(p printer, state State, err error) error
}

var DefaultDumper Dumper = &ConsoleDumper{}
