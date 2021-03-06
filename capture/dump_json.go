package capture

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/podhmo/tenuki/capture/httputil"
)

type JsonDumper struct {
}

func (d *JsonDumper) DumpRequest(p printer, req *http.Request) (State, error) {
	info, err := httputil.DumpRequestJSON(req, true /* body */)
	if err != nil {
		return nil, err
	}

	// TODO: use printer
	enc := json.NewEncoder(os.Stderr)
	enc.SetIndent("", "  ")
	return nil, enc.Encode(info)
}
func (d *JsonDumper) DumpError(p printer, state State, err error) error {
	p.Printf("\x1b[90merror:\n%+v\x1b[0m", err)
	return err
}

func (d *JsonDumper) DumpResponse(p printer, state State, res *http.Response) error {
	info, err := httputil.DumpResponseJSON(res, true /* body */)
	if err != nil {
		return err
	}

	// TODO: use printer
	enc := json.NewEncoder(os.Stderr)
	enc.SetIndent("", "  ")
	return enc.Encode(info)
}

var _ Dumper = &JsonDumper{}
