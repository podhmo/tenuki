package openapistyle

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type RequestInfo struct {
	Method      string `json:"method"`
	URL         string `json:"url"`
	HTTPVersion string `json:"httpVersion,omitempty"`
	HeaderSize  int64  `json:"headerSize,omitempty" default:"-1"`
	BodySize    int64  `json:"bodySize,omitempty" default:"-1"`

	Paths Paths `json:"paths"`

	body io.ReadCloser `json:"-"` // for interface
}

func (info RequestInfo) Body() io.ReadCloser {
	return info.body
}

func ExtractRequestInfo(req *http.Request, body io.ReadCloser) (RequestInfo, error) {
	info := RequestInfo{body: body}

	{
		reqURI := req.RequestURI
		if reqURI == "" {
			reqURI = req.URL.RequestURI()
		}

		absRequestURI := strings.HasPrefix(req.RequestURI, "http://") || strings.HasPrefix(req.RequestURI, "https://")
		if !absRequestURI {
			reqURI = fmt.Sprintf("%s://%s%s", valueOrDefault(req.URL.Scheme, "https"), req.URL.Host, reqURI)
		}
		info.URL = reqURI
	}

	info.Method = valueOrDefault(req.Method, "GET")
	info.HTTPVersion = req.Proto
	info.HeaderSize = -1 // TODO
	info.BodySize = -1   // TODO

	paths, err := toPaths(req, body)
	if err != nil {
		return info, fmt.Errorf("extract paths, %w", err)
	}
	info.Paths = paths
	return info, nil
}

type ResponseInfo struct {
	body io.ReadCloser `json:"-"` // for interface
}

func (info ResponseInfo) Body() io.ReadCloser {
	return info.body
}

func ExtractResponseInfo(resp *http.Response) ResponseInfo {
	info := ResponseInfo{body: resp.Body}
	return info
}

// Return value if nonempty, def otherwise.
func valueOrDefault(value, def string) string {
	if value != "" {
		return value
	}
	return def
}
