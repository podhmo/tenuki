package capture

import (
	"net/http"

	"github.com/podhmo/tenuki/capture/httputil"
)

type ConsoleDumper struct {
}

func (d *ConsoleDumper) DumpRequest(p printer, req *http.Request) (State, error) {
	b, err := httputil.DumpRequest(req, true /* body */)
	if err != nil {
		return nil, err
	}

	p.Printf("\x1b[90mrequest:\n%s\x1b[0m", string(b))
	return nil, nil
}
func (d *ConsoleDumper) DumpError(p printer, state State, err error) error {
	p.Printf("\x1b[90merror:\n%+v\x1b[0m", err)
	return err
}

func (d *ConsoleDumper) DumpResponse(p printer, state State, res *http.Response) error {
	b, err := httputil.DumpResponse(res, true /* body */)
	if err != nil {
		return err
	}

	p.Printf("\x1b[90mresponse:\n%s\x1b[0m", string(b))
	return nil
}

var _ Dumper = &ConsoleDumper{}
