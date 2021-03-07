package capture

import (
	"github.com/podhmo/tenuki/capture/gostyle"
	"github.com/podhmo/tenuki/capture/httputil"
	"github.com/podhmo/tenuki/capture/openapistyle"
)

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
