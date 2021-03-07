package capture

import (
	"io"
	"net/http"
)

// for text output

type HTTPutilDumpRequestFunc func(req *http.Request, body bool) ([]byte, error)

func (f HTTPutilDumpRequestFunc) Extract(req *http.Request) ([]byte, error) {
	return f(req, true /* body */)
}

type HTTPutilDumpResponseFunc func(resp *http.Response, body bool) ([]byte, error)

func (f HTTPutilDumpResponseFunc) Extract(resp *http.Response) ([]byte, error) {
	return f(resp, true /* body */)
}

// for json output
type JSONDumpRequestFuncWithStyle struct {
	Dump func(
		req *http.Request,
		body bool,
		extractInfo func(*http.Request, io.Reader) (interface{ Info() interface{} }, error)) ([]byte, error)
	Style func(
		*http.Request,
		io.Reader,
	) (interface{ Info() interface{} }, error)
}

func (f *JSONDumpRequestFuncWithStyle) Extract(req *http.Request) ([]byte, error) {
	return f.Dump(req, true /* body */, f.Style)
}

type JSONDumpResponseFuncWithStyle struct {
	Dump func(
		req *http.Response,
		body bool,
		extractInfo func(*http.Response, io.Reader) (interface{ Info() interface{} }, error)) ([]byte, error)
	Style func(
		*http.Response,
		io.Reader,
	) (interface{ Info() interface{} }, error)
}

func (f *JSONDumpResponseFuncWithStyle) Extract(req *http.Response) ([]byte, error) {
	return f.Dump(req, true /* body */, f.Style)
}
