package capture

import (
	"io"
	"net/http"

	"github.com/podhmo/tenuki/capture/httputil"
)

type Layout struct {
	Request interface {
		Extract(*http.Request) (State, error)
	}
	Response interface {
		Extract(*http.Response) (State, error)
	}
}

type HTTPutilDumpRequestFunc func(req *http.Request, body bool) ([]byte, error)

func (f HTTPutilDumpRequestFunc) Extract(req *http.Request) (State, error) {
	b, err := f(req, true /* body */)
	if err != nil {
		return nil, err
	}
	return bytesLazy(b), nil
}

type HTTPutilDumpResponseFunc func(resp *http.Response, body bool) ([]byte, error)

func (f HTTPutilDumpResponseFunc) Extract(resp *http.Response) (State, error) {
	b, err := f(resp, true /* body */)
	if err != nil {
		return nil, err
	}
	return bytesLazy(b), nil
}

// for json output
type JSONDumpRequestFuncWithStyle struct {
	Dump func(
		req *http.Request,
		body bool,
		extractInfo func(*http.Request, io.Reader) (interface{ Info() interface{} }, error)) (*httputil.JSONState, error)
	Style func(
		*http.Request,
		io.Reader,
	) (interface{ Info() interface{} }, error)
}

func (f *JSONDumpRequestFuncWithStyle) Extract(req *http.Request) (State, error) {
	return f.Dump(req, true /* body */, f.Style)
}

type JSONDumpResponseFuncWithStyle struct {
	Dump func(
		req *http.Response,
		body bool,
		extractInfo func(*http.Response, io.Reader) (interface{ Info() interface{} }, error)) (*httputil.JSONState, error)
	Style func(
		*http.Response,
		io.Reader,
	) (interface{ Info() interface{} }, error)
}

func (f *JSONDumpResponseFuncWithStyle) Extract(req *http.Response) (State, error) {
	return f.Dump(req, true /* body */, f.Style)
}
