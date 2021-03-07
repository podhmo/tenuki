package capture

import (
	"net/http"
	"unsafe"
)

type ConsoleTransport struct {
	Transport http.RoundTripper
	Layout    *Layout
	Printer   interface {
		Printf(fmt string, args ...interface{})
	}
}

func (ct *ConsoleTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := ct.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	s, err := ct.DumpRequest(req)
	if err != nil {
		return nil, err
	}
	res, err := transport.RoundTrip(req)
	if err != nil {
		return nil, ct.DumpError(err)
	}
	if err := ct.DumpResponse(res, s); err != nil {
		return nil, err
	}
	return res, nil
}

func (ct *ConsoleTransport) DumpRequest(req *http.Request) (State, error) {
	layout := ct.Layout
	if layout == nil {
		layout = DefaultLayout
	}
	s, err := layout.Request.Extract(req)
	if err != nil {
		return nil, err
	}
	b, err := s.Encode()
	if err != nil {
		return nil, err
	}
	ct.Printer.Printf("\x1b[90mrequest:\n%s\x1b[0m", *(*string)(unsafe.Pointer(&b)))
	return s, nil
}

func (ct *ConsoleTransport) DumpError(err error) error {
	ct.Printer.Printf("\x1b[90merror:\n%+v\x1b[0m", err)
	return err
}

func (ct *ConsoleTransport) DumpResponse(res *http.Response, s State) error {
	layout := ct.Layout
	if layout == nil {
		layout = DefaultLayout
	}
	b, err := layout.Response.Extract(res, s)
	if err != nil {
		return err
	}

	ct.Printer.Printf("\x1b[90mresponse:\n%s\x1b[0m", *(*string)(unsafe.Pointer(&b)))
	return nil
}
