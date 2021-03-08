package openapistyle

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/podhmo/tenuki/capture/style"
)

type RequestInfo struct {
	Method      string `json:"method"`
	URL         string `json:"url"`
	HTTPVersion string `json:"httpVersion,omitempty"`
	HeaderSize  int64  `json:"headerSize,omitempty" default:"-1"`
	BodySize    int64  `json:"bodySize,omitempty" default:"-1"`

	Paths Paths `json:"paths"`
}

func (s *RequestInfo) Merge(res style.Info) style.Info {
	// TODO
	return res
}

// TODO
func (s *RequestInfo) HandleError(open func() (io.WriteCloser, error), err error) {
}

func ExtractRequestInfo(req *http.Request) (style.Info, error) {
	info := RequestInfo{}
	body := req.Body

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
		return nil, fmt.Errorf("extract paths, %w", err)
	}
	info.Paths = paths
	return &info, nil
}

func ExtractResponseInfo(resp *http.Response) (style.Info, error) {
	// TODO:
	return &RequestInfo{}, nil
}

// Return value if nonempty, def otherwise.
func valueOrDefault(value, def string) string {
	if value != "" {
		return value
	}
	return def
}
