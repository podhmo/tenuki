package capture

import (
	"net/http"

	"github.com/podhmo/tenuki/capture/gostyle"
	"github.com/podhmo/tenuki/capture/httputil"
	"github.com/podhmo/tenuki/capture/openapistyle"
)

type CapturedTransport struct {
}

func (ct *CapturedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := ct.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	dumper := ct.Dumper
	if dumper == nil {
		dumper = DefaultDumper
	}

	s, err := dumper.DumpRequest(ct.Printer, req)
	if err != nil {
		return nil, err
	}
	res, err := transport.RoundTrip(req)
	if err != nil {
		return nil, dumper.DumpError(ct.Printer, s, err)
	}
	if err := dumper.DumpResponse(ct.Printer, s, res); err != nil {
		return nil, err
	}
	return res, nil
}

type State interface {
	Request() *http.Request
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
