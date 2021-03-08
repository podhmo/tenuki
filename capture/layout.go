package capture

import (
	"net/http"

	"github.com/podhmo/tenuki/capture/httputil"
	"github.com/podhmo/tenuki/capture/style"
)

type Layout struct {
	Request interface {
		Extract(*http.Request) (style.State, error)
	}
	Response interface {
		Extract(*http.Response) (style.State, error)
	}
}

type HTTPutilDumpRequestFunc func(req *http.Request, body bool) ([]byte, error)

func (f HTTPutilDumpRequestFunc) Extract(req *http.Request) (style.State, error) {
	b, err := f(req, true /* body */)
	if err != nil {
		return nil, err
	}
	return &bytesState{b: b, req: req}, nil
}

type HTTPutilDumpResponseFunc func(resp *http.Response, body bool) ([]byte, error)

func (f HTTPutilDumpResponseFunc) Extract(resp *http.Response) (style.State, error) {
	b, err := f(resp, true /* body */)
	if err != nil {
		return nil, err
	}
	return &bytesState{b: b}, nil
}

// for json output
type JSONDumpRequestFuncWithStyle struct {
	Style func(*http.Request) (style.Info, error)
}

func (f *JSONDumpRequestFuncWithStyle) Extract(req *http.Request) (style.State, error) {
	info, err := httputil.DumpRequestJSON(req, true /* body */, f.Style)
	return &jsonState{info: info}, err
}

type JSONDumpResponseFuncWithStyle struct {
	Style func(*http.Response) (style.Info, error)
}

func (f *JSONDumpResponseFuncWithStyle) Extract(res *http.Response) (style.State, error) {
	info, err := httputil.DumpResponseJSON(res, true /* body */, f.Style)
	return &jsonState{info: info}, err
}
