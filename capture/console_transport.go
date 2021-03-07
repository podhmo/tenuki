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
	if err := ct.DumpRequest(req); err != nil {
		return nil, err
	}
	res, err := transport.RoundTrip(req)
	if err != nil {
		return nil, ct.DumpError(err)
	}
	if err := ct.DumpResponse(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (ct *ConsoleTransport) DumpRequest(req *http.Request) error {
	layout := ct.Layout
	if layout == nil {
		layout = DefaultLayout
	}
	b, err := layout.Request.Extract(req)
	if err != nil {
		return err
	}

	ct.Printer.Printf("\x1b[90mrequest:\n%s\x1b[0m", *(*string)(unsafe.Pointer(&b)))
	return nil
}

func (ct *ConsoleTransport) DumpError(err error) error {
	ct.Printer.Printf("\x1b[90merror:\n%+v\x1b[0m", err)
	return err
}

func (ct *ConsoleTransport) DumpResponse(res *http.Response) error {
	layout := ct.Layout
	if layout == nil {
		layout = DefaultLayout
	}
	b, err := layout.Response.Extract(res)
	if err != nil {
		return err
	}

	ct.Printer.Printf("\x1b[90mresponse:\n%s\x1b[0m", *(*string)(unsafe.Pointer(&b)))
	return nil
}
