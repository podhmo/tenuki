package capture

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/podhmo/tenuki/capture/gostyle"
	"github.com/podhmo/tenuki/capture/httputil"
)

type JSONDumper struct {
	ExtractRequestInfo  func(*http.Request, io.Reader) (interface{ Info() interface{} }, error)
	ExtractResponseInfo func(*http.Response, io.Reader) (interface{ Info() interface{} }, error)
}

func (d *JSONDumper) DumpRequest(p printer, req *http.Request) (State, error) {
	extractInfo := gostyle.ExtractRequestInfo
	if d.ExtractRequestInfo != nil {
		extractInfo = d.ExtractRequestInfo
	}
	info, err := httputil.DumpRequestJSON(req, true /* body */, extractInfo)
	if err != nil {
		return nil, err
	}

	// TODO: use printer
	enc := json.NewEncoder(os.Stderr)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return nil, enc.Encode(info)
}
func (d *JSONDumper) DumpError(p printer, state State, err error) error {
	p.Printf("\x1b[90merror:\n%+v\x1b[0m", err)
	return err
}

func (d *JSONDumper) DumpResponse(p printer, state State, res *http.Response) error {
	extractInfo := gostyle.ExtractResponseInfo
	if d.ExtractResponseInfo != nil {
		extractInfo = d.ExtractResponseInfo
	}

	info, err := httputil.DumpResponseJSON(res, true /* body */, extractInfo)
	if err != nil {
		return err
	}

	// TODO: use printer
	enc := json.NewEncoder(os.Stderr)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(info)
}

var _ Dumper = &JSONDumper{}
