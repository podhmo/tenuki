package capture

import (
	"io"
	"net/http"

	"github.com/podhmo/tenuki/capture/gostyle"
	"github.com/podhmo/tenuki/capture/httputil"
	"github.com/podhmo/tenuki/capture/openapistyle"
)

type Layout struct {
	Request interface {
		Extract(*http.Request) (State, error)
	}
	Response interface {
		Extract(*http.Response) (State, error)
	}
}

type State interface {
	Encode() ([]byte, error)
	Emit(w io.WriteCloser) error
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
	DefaultLayout = TextLayout
)
