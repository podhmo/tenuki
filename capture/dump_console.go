package capture

import (
	"net/http"
	"unsafe"
)

type ConsoleDumper struct {
	Layout *Layout
}

func (d *ConsoleDumper) DumpRequest(p printer, req *http.Request) (State, error) {
	layout := d.Layout
	if layout == nil {
		layout = DefaultLayout
	}
	b, err := layout.Request.Extract(req)
	if err != nil {
		return nil, err
	}

	p.Printf("\x1b[90mrequest:\n%s\x1b[0m", *(*string)(unsafe.Pointer(&b)))
	return nil, nil
}
func (d *ConsoleDumper) DumpError(p printer, state State, err error) error {
	p.Printf("\x1b[90merror:\n%+v\x1b[0m", err)
	return err
}

func (d *ConsoleDumper) DumpResponse(p printer, state State, res *http.Response) error {
	layout := d.Layout
	if layout == nil {
		layout = DefaultLayout
	}
	b, err := layout.Response.Extract(res)
	if err != nil {
		return err
	}

	p.Printf("\x1b[90mresponse:\n%s\x1b[0m", *(*string)(unsafe.Pointer(&b)))
	return nil
}

var _ Dumper = &ConsoleDumper{}
