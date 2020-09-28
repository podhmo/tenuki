package capture

import (
	"net/http"
	"net/http/httputil"
)

type ConsoleDumper struct {
}

func (d *ConsoleDumper) DumpRequest(p printer, req *http.Request) error {
	b, err := httputil.DumpRequest(req, true /* body */)
	if err != nil {
		return err
	}

	p.Printf("\x1b[90mrequest:\n%s\x1b[0m", string(b))
	return nil
}
func (d *ConsoleDumper) DumpResponse(p printer, res *http.Response) error {
	b, err := httputil.DumpResponse(res, true /* body */)
	if err != nil {
		return err
	}

	p.Printf("\x1b[90mresponse:\n%s\x1b[0m", string(b))
	return nil
}
