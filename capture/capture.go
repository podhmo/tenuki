package capture

import (
	"net/http"

	"github.com/podhmo/tenuki/capture/gostyle"
	"github.com/podhmo/tenuki/capture/httputil"
	"github.com/podhmo/tenuki/capture/openapistyle"
)

type printer interface {
	Printf(fmt string, args ...interface{})
}

type State interface {
	Request() *http.Request
}
type Dumper interface {
	DumpRequest(p printer, req *http.Request) (State, error)
	DumpResponse(p printer, state State, res *http.Response) error
	DumpError(p printer, state State, err error) error
}

type Layout struct {
	Request interface {
		Extract(*http.Request) ([]byte, error)
	}
	Response interface {
		Extract(*http.Response) ([]byte, error)
	}
}

var (
	TextLayout = &Layout{
		Request:  HTTPutilDumpRequestFunc(httputil.DumpRequestOut),
		Response: HTTPutilDumpResponseFunc(httputil.DumpResponse),
	}
	JSONLayout = &Layout{
		Request: &JSONDumpRequestFuncWithStyle{
			Dump:  httputil.DumpRequestJSON,
			Style: gostyle.ExtractRequestInfo,
		},
		Response: &JSONDumpResponseFuncWithStyle{
			Dump:  httputil.DumpResponseJSON,
			Style: gostyle.ExtractResponseInfo,
		},
	}
	OpenAPILayout = &Layout{
		Request: &JSONDumpRequestFuncWithStyle{
			Dump:  httputil.DumpRequestJSON,
			Style: openapistyle.ExtractRequestInfo,
		},
		Response: &JSONDumpResponseFuncWithStyle{
			Dump:  httputil.DumpResponseJSON,
			Style: openapistyle.ExtractResponseInfo,
		},
	}
)

// default setting
var (
	DefaultDumper Dumper = &ConsoleDumper{}
	DefaultLayout        = TextLayout
)
